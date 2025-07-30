package rpc

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/features/handler"
	registrar_service_v1 "github.com/fadliarz/distributed-faas/services/registrar-service/gen/go/registrar-service/v1"
	"google.golang.org/grpc"
)

type RegistrarServer struct {
	registrar_service_v1.UnimplementedRegistrarServiceServer
	handler *handler.CommandHandler
}

func NewRegistrarServer(handler *handler.CommandHandler) *RegistrarServer {
	return &RegistrarServer{
		handler: handler,
	}
}

func (s *RegistrarServer) Register(server *grpc.Server) {
	registrar_service_v1.RegisterRegistrarServiceServer(server, s)
}

func (s *RegistrarServer) RegisterMachine(ctx context.Context,
	req *registrar_service_v1.RegisterMachineRequest,
) (*registrar_service_v1.RegisterMachineResponse, error) {
	cmd := &command.CreateMachineCommand{
		Address: req.Address,
	}

	machineID, err := s.handler.CreateMachine(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &registrar_service_v1.RegisterMachineResponse{
		MachineId: machineID.String(),
		Status:    "success",
		Message:   "Machine registered successfully",
	}, nil
}
