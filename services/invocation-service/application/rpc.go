package application

import (
	"context"
	"errors"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	invocation_service_v1 "github.com/fadliarz/distributed-faas/services/invocation-service/gen/go/invocation-service/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InvocationServer struct {
	invocation_service_v1.UnimplementedInvocationServiceServer

	handler *application.CommandHandler
}

func NewInvocationServer(handler *application.CommandHandler) (*InvocationServer, error) {
	return &InvocationServer{handler: handler}, nil
}

func (s *InvocationServer) Register(server *grpc.Server) {
	invocation_service_v1.RegisterInvocationServiceServer(server, s)
}

func (s *InvocationServer) CreateInvocation(ctx context.Context, req *invocation_service_v1.CreateInvocationRequest) (*invocation_service_v1.CreateInvocationResponse, error) {
	cmd := &application.CreateInvocationCommand{
		UserID:     req.UserId,
		FunctionID: req.FunctionId,
	}

	log.Debug().Msgf("Creating invocation for user %s and function %s", cmd.UserID, cmd.FunctionID)

	invocationID, err := s.handler.CreateInvocation(context.Background(), cmd)

	if err != nil && errors.Is(err, domain.ErrUserNotAuthorized) {
		log.Debug().Err(err).Msgf("User %s is not authorized to create invocation for function %s", cmd.UserID, cmd.FunctionID)

		return nil, status.Error(codes.PermissionDenied, "you are not authorized to perform this action")
	} else if err != nil {
		log.Debug().Err(err).Msgf("Failed to create invocation for user %s and function %s", cmd.UserID, cmd.FunctionID)

		return nil, status.Error(codes.Internal, "failed to create invocation")
	}

	return &invocation_service_v1.CreateInvocationResponse{
		InvocationId: invocationID.String(),
		Status:       "success",
		Message:      "Invocation created successfully",
	}, nil
}
