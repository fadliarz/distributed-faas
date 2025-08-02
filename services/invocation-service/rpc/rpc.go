package rpc

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

func NewInvocationServer(handler *application.CommandHandler) *InvocationServer {
	return &InvocationServer{
		handler: handler,
	}
}

func (s *InvocationServer) Register(server *grpc.Server) {
	invocation_service_v1.RegisterInvocationServiceServer(server, s)
}

func (s *InvocationServer) CreateInvocation(ctx context.Context, req *invocation_service_v1.CreateInvocationRequest) (*invocation_service_v1.CreateInvocationResponse, error) {
	cmd := &application.CreateInvocationCommand{
		UserID:     req.UserId,
		FunctionID: req.FunctionId,
	}

	log.Info().Msgf("Creating invocation for user %s and function %s", req.UserId, req.FunctionId)

	invocation, err := s.handler.CreateInvocation(ctx, cmd)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create invocation")

		if errors.Is(err, domain.ErrUserNotAuthorized) {
			return nil, status.Errorf(codes.PermissionDenied, "User not authorized to invoke this function")
		}

		return nil, status.Errorf(codes.Internal, "Failed to create invocation: %v", err)
	}

	return &invocation_service_v1.CreateInvocationResponse{
		InvocationId:  invocation.InvocationID.String(),
		FunctionId:    invocation.FunctionID.String(),
		UserId:        invocation.UserID.String(),
		SourceCodeUrl: invocation.SourceCodeURL.String(),
		OutputUrl:     invocation.OutputURL.String(),
		Status:        "success",
		Message:       "Invocation created successfully",
	}, nil
}

func (s *InvocationServer) GetInvocation(ctx context.Context, req *invocation_service_v1.GetInvocationRequest) (*invocation_service_v1.GetInvocationResponse, error) {
	log.Info().Msgf("Fetching invocation with ID %s", req.InvocationId)

	invocation, err := s.handler.GetInvocation(ctx, &application.GetInvocationQuery{
		InvocationID: req.InvocationId,
		UserID:       req.UserId,
	})

	if err != nil {
		log.Debug().Err(err).Msgf("Failed to fetch invocation with ID %s", req.InvocationId)

		if errors.Is(err, domain.ErrUserNotAuthorized) {
			return nil, status.Errorf(codes.PermissionDenied, "User not authorized to access this invocation")
		}

		return nil, status.Errorf(codes.Internal, "Failed to fetch invocation")
	}

	return &invocation_service_v1.GetInvocationResponse{
		InvocationId:  invocation.InvocationID.String(),
		FunctionId:    invocation.FunctionID.String(),
		SourceCodeUrl: invocation.SourceCodeURL.String(),
		OutputUrl:     invocation.OutputURL.String(),
	}, nil
}
