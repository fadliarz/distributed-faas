package rpc

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/machine/domain/application-service"
	machine_service_v1 "github.com/fadliarz/distributed-faas/services/machine/gen/go/machine-service/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type MachineServer struct {
	machine_service_v1.UnimplementedMachineServiceServer

	ctx     context.Context
	handler *application.CommandHandler
}

func NewMachineServer(ctx context.Context, handler *application.CommandHandler) *MachineServer {
	return &MachineServer{
		ctx:     ctx,
		handler: handler,
	}
}

func (s *MachineServer) Register(server *grpc.Server) {
	machine_service_v1.RegisterMachineServiceServer(server, s)
}

func (s *MachineServer) ExecuteFunction(ctx context.Context, req *machine_service_v1.ExecuteFunctionRequest) (*machine_service_v1.ExecuteFunctionResponse, error) {
	cmd := &application.ProcessInvocationCommand{
		InvocationID:  req.InvocationId,
		FunctionID:    req.FunctionId,
		UserID:        req.UserId,
		SourceCodeURL: req.SourceCodeUrl,
		Status:        req.Status,
		Timestamp:     req.Timestamp,
	}

	log.Debug().Msgf("Received ExecuteFunction request: %+v", cmd)

	_, err := s.handler.ProcessInvocation(s.ctx, cmd)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to process invocation")

		return nil, fmt.Errorf("Failed to process invocation")
	}

	return &machine_service_v1.ExecuteFunctionResponse{
		Status:  "SUCCESS",
		Message: "Function execution request received and logged",
	}, nil
}
