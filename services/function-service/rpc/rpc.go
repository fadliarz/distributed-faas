package rpc

import (
	"context"
	"errors"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service"
	function_service_v1 "github.com/fadliarz/distributed-faas/services/function-service/gen/go/function-service/v1"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FunctionServer struct {
	function_service_v1.UnimplementedFunctionServiceServer

	handler *application.CommandHandler
}

func NewFunctionServer(handler *application.CommandHandler) *FunctionServer {
	return &FunctionServer{
		handler: handler,
	}
}

func (s *FunctionServer) Register(server *grpc.Server) {
	function_service_v1.RegisterFunctionServiceServer(server, s)
}

func (s *FunctionServer) CreateFunction(ctx context.Context, req *function_service_v1.CreateFunctionRequest) (*function_service_v1.CreateFunctionResponse, error) {
	cmd := &application.CreateFunctionCommand{
		UserID:        req.UserId,
		SourceCodeURL: req.SourceCodeUrl,
	}

	functionID, err := s.handler.CreateFunction(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return &function_service_v1.CreateFunctionResponse{
		FunctionId: functionID.String(),
		Status:     "success",
		Message:    "Function created successfully",
	}, nil
}

func (s *FunctionServer) GetFunctionUploadPresignedURL(ctx context.Context, req *function_service_v1.GetFunctionUploadPresignedURLRequest) (*function_service_v1.GetFunctionUploadPresignedURLResponse, error) {
	query := &application.GetFunctionUploadPresignedURLQuery{
		UserID:     req.UserId,
		FunctionID: req.FunctionId,
	}

	presignedURL, err := s.handler.GetFunctionUploadPresignedURL(ctx, query)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotAuthorized) {
			return nil, status.Errorf(codes.PermissionDenied, "user not authorized: %v", err)
		}

		return nil, status.Errorf(codes.Internal, "failed to get presigned URL: %v", err)
	}

	return &function_service_v1.GetFunctionUploadPresignedURLResponse{
		PresignedUrl: presignedURL,
		Status:       "success",
		Message:      "Presigned URL retrieved successfully",
	}, nil
}
