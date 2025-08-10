package rpc

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/user-service/domain/application-service"
	user_service_v1 "github.com/fadliarz/distributed-faas/services/user-service/gen/go/user-service/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	user_service_v1.UnimplementedUserServiceServer

	handler *application.UserCommandHandler
}

func NewUserServer(handler *application.UserCommandHandler) *UserServer {
	return &UserServer{
		handler: handler,
	}
}

func (s *UserServer) Register(server *grpc.Server) {
	user_service_v1.RegisterUserServiceServer(server, s)
}

func (s *UserServer) CreateUser(ctx context.Context, req *user_service_v1.CreateUserRequest) (*user_service_v1.CreateUserResponse, error) {
	command := &application.CreateUserCommand{
		Password: req.Password,
	}

	user, err := s.handler.CreateUser(ctx, command)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create user")

		return nil, status.Error(codes.Internal, "Failed to create user")
	}

	return &user_service_v1.CreateUserResponse{
		UserId: user.UserID.String(),
	}, nil
}
