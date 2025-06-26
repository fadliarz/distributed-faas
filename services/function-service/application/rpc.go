package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/features/handler"
	function_service_v1 "github.com/fadliarz/distributed-faas/services/function-service/gen/go/function-service/v1"
	"google.golang.org/grpc"
)

type FunctionServer struct {
	function_service_v1.UnimplementedFunctionServiceServer
	handler *handler.CommandHandler
}

func NewFunctionServer() *FunctionServer {
	return &FunctionServer{handler: handler.NewCommandHandler()}
}

func (s *FunctionServer) Register(server *grpc.Server) {
	function_service_v1.RegisterFunctionServiceServer(server, s)
}

func (s *FunctionServer) CreateFunction(ctx context.Context, req *function_service_v1.CreateFunctionRequest) (*function_service_v1.CreateFunctionResponse, error) {
	cmd := &command.CreateFunctionCommand{
		UserID:        req.UserId,
		SourceCodeURL: req.SourceCodeUrl,
	}

	functionID, err := s.handler.CreateFunction(cmd)
	if err != nil {
		return nil, err
	}

	return &function_service_v1.CreateFunctionResponse{
		FunctionId: functionID.String(),
		Status:     "success",
		Message:    "Function created successfully",
	}, nil
}
