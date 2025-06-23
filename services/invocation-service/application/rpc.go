package application

import (
	"context"

	"github.com/fadliarz/services/invocation-service/domain/application-service/features/command"
	"github.com/fadliarz/services/invocation-service/domain/application-service/features/handler"
	"github.com/fadliarz/services/invocation-service/domain/domain-core/core"
	invocation_service_v1 "github.com/fadliarz/services/invocation-service/gen/go/invocation-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InvocationServer struct {
	invocation_service_v1.UnimplementedInvocationServiceServer
	handler *handler.CommandHandler
}

func NewInvocationServer() *InvocationServer {
	return &InvocationServer{handler: handler.NewCommandHandler()}
}

func (s *InvocationServer) Register(server *grpc.Server) {
	invocation_service_v1.RegisterInvocationServiceServer(server, s)
}

func (s *InvocationServer) CreateInvocation(ctx context.Context, req *invocation_service_v1.CreateInvocationRequest) (*invocation_service_v1.CreateInvocationResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be empty")
	}

	if req.FunctionId == "" {
		return nil, status.Error(codes.InvalidArgument, "function_id is required")
	}

	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	cmd := &command.CreateInvocationCommand{
		UserID:     req.UserId,
		FunctionID: req.FunctionId,
	}

	invocationID, err := s.handler.CreateInvocation(cmd)
	if err != nil {
		if core.IsErrorType(err, core.ValidationError) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		} else if core.IsErrorType(err, core.NotFoundError) {
			return nil, status.Error(codes.NotFound, err.Error())
		} else if core.IsErrorType(err, core.DatabaseError) {
			return nil, status.Error(codes.Internal, "database operation failed")
		} else {
			return nil, status.Error(codes.Internal, "internal server error")
		}
	}

	return &invocation_service_v1.CreateInvocationResponse{
		InvocationId: invocationID.String(),
		Status:       "success",
		Message:      "Invocation created successfully",
	}, nil
}
