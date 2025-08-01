package rpc

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service"
	registrar_service_v1 "github.com/fadliarz/distributed-faas/services/registrar-service/gen/go/registrar-service/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RegistrarServer struct {
	registrar_service_v1.UnimplementedRegistrarServiceServer

	handler *application.CommandHandler
}

func NewRegistrarServer(handler *application.CommandHandler) *RegistrarServer {
	return &RegistrarServer{
		handler: handler,
	}
}

func (s *RegistrarServer) Register(server *grpc.Server) {
	registrar_service_v1.RegisterRegistrarServiceServer(server, s)
}

func (s *RegistrarServer) RegisterMachine(ctx context.Context, req *registrar_service_v1.RegisterMachineRequest) (*registrar_service_v1.RegisterMachineResponse, error) {
	cmd := &application.CreateMachineCommand{
		Address: req.Address,
	}

	machine, err := s.handler.CreateMachine(ctx, cmd)
	if err != nil {
		log.Warn().Err(err).Msg("failed to create machine")

		return nil, status.Errorf(codes.Internal, "failed to register machine")
	}

	return &registrar_service_v1.RegisterMachineResponse{
		MachineId: machine.MachineID.String(),
		Address:   machine.Address.String(),
		Status:    "success",
		Message:   "Machine registered successfully",
	}, nil
}
