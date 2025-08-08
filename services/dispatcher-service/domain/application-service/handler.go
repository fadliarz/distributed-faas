package application

import (
	"context"
	"fmt"
	"time"

	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/domain-core"
	machine_service_v1 "github.com/fadliarz/distributed-faas/services/dispatcher-service/gen/go/machine-service/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Constructors

type InvocationEventHandler struct {
	repositoryManager *InvocationEventHandlerRepositoryManager
}

type InvocationEventHandlerRepositoryManager struct {
	Machine MachineRepository
}

func NewInvocationEventHandler(repositoryManager *InvocationEventHandlerRepositoryManager) *InvocationEventHandler {
	return &InvocationEventHandler{
		repositoryManager: repositoryManager,
	}
}

func NewInvocationEventHandlerRepositoryManager(machine MachineRepository) *InvocationEventHandlerRepositoryManager {
	return &InvocationEventHandlerRepositoryManager{
		Machine: machine,
	}
}

// Methods

func (h *InvocationEventHandler) HandleInvocationCreatedEvent(ctx context.Context, event *InvocationCreatedEvent) error {
	for true {
		var machines []domain.Machine
		var err error

		for machines, err = h.repositoryManager.Machine.FindManyByStatus(ctx, domain.NewStatusFromInt(int(domain.Available))); err != nil || len(machines) == 0; {
			if err != nil {
				log.Warn().Err(err).Msg("Failed to find machines by status, retrying...")

				continue
			}

			if len(machines) > 0 {
				log.Warn().Msgf("Found %d available machines, but still waiting for more to become available", len(machines))
			}

			log.Info().Msg("Pausing for 5 seconds before retrying to find available machines")

			time.Sleep(5 * time.Second)
		}

		for i := 0; i < len(machines); i++ {
			err = h.executeInvocation(ctx, machines[0], event)
			if err == nil {
				log.Info().Msgf("Invocation processed successfully on machine %s", machines[0].MachineID.String())

				return nil
			}

			log.Warn().Err(err).Msgf("Failed to process invocation on machine %s, trying next machine", machines[0].MachineID.String())
		}

		log.Warn().Msg("No machines successfully processed the invocation, retrying in 5 seconds")

		time.Sleep(5 * time.Second)
	}

	return nil
}

func (h *InvocationEventHandler) executeInvocation(ctx context.Context, machine domain.Machine, event *InvocationCreatedEvent) error {
	// Create a gRPC client
	address := machine.Address.String()
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return fmt.Errorf("failed to connect to machine service at %s: %w", address, err)
	}
	defer conn.Close()

	client := machine_service_v1.NewMachineServiceClient(conn)

	// Create the gRPC request and send the gRPC request
	request := &machine_service_v1.ExecuteFunctionRequest{
		InvocationId:  event.InvocationID,
		FunctionId:    event.FunctionID,
		UserId:        event.UserID,
		SourceCodeUrl: event.SourceCodeURL,
		Status:        event.Status,
		Timestamp:     event.Timestamp,
	}

	_, err = client.ExecuteFunction(ctx, request)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to execute function on machine %s", machine.MachineID.String())

		return fmt.Errorf("failed to execute function on machine %s: %w", machine.MachineID.String(), err)
	}

	return nil
}
