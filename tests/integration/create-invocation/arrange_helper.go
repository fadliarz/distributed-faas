package main

import (
	"testing"

	registrar_service_v1 "github.com/fadliarz/distributed-faas/services/registrar-service/gen/go/registrar-service/v1"
	function_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/function-service/v1"
	invocation_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/invocation-service/v1"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ArrangeHelper struct {
	t      *testing.T
	config *TestConfig
}

func NewArrangeHelper(t *testing.T, config *TestConfig) *ArrangeHelper {
	return &ArrangeHelper{
		t:      t,
		config: config,
	}
}

func (h *ArrangeHelper) CreateFunction() *function_service_v1.CreateFunctionResponse {
	conn, err := grpc.NewClient(h.config.GrpcEndpoints.FunctionService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := function_service_v1.NewFunctionServiceClient(conn)

	req := function_service_v1.CreateFunctionRequest{
		UserId: uuid.NewString(),
	}

	response, err := client.CreateFunction(h.t.Context(), &req)

	require.NoError(h.t, err, "Failed to create function")

	require.NotEmpty(h.t, response.GetFunctionId(), "Function ID should not be empty")
	require.Equal(h.t, req.UserId, response.GetUserId(), "User ID should match the request")
	require.Empty(h.t, response.GetSourceCodeUrl(), "Source code URL should be empty")

	return response
}

func (h *ArrangeHelper) UpdateFunctionSourceCodeURL(userID string, functionID string, sourceCodeURL string) *function_service_v1.UpdateFunctionSourceCodeURLResponse {
	conn, err := grpc.NewClient(h.config.GrpcEndpoints.FunctionService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := function_service_v1.NewFunctionServiceClient(conn)

	req := function_service_v1.UpdateFunctionSourceCodeURLRequest{
		UserId:        userID,
		FunctionId:    functionID,
		SourceCodeUrl: sourceCodeURL,
	}

	response, err := client.UpdateFunctionSourceCodeURL(h.t.Context(), &req)

	require.NoError(h.t, err, "Failed to update function source code URL")

	return response
}

func (h *ArrangeHelper) CreateInvocation(userID string, functionID string) (*invocation_service_v1.CreateInvocationResponse, error) {
	conn, err := grpc.NewClient(h.config.GrpcEndpoints.InvocationService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := invocation_service_v1.NewInvocationServiceClient(conn)

	req := invocation_service_v1.CreateInvocationRequest{
		UserId:     userID,
		FunctionId: functionID,
	}

	response, err := client.CreateInvocation(h.t.Context(), &req)


	return response, err
}

func (h *ArrangeHelper) RegisterMachine() *registrar_service_v1.RegisterMachineResponse {
	conn, err := grpc.NewClient(h.config.GrpcEndpoints.RegistrarService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := registrar_service_v1.NewRegistrarServiceClient(conn)

	req := registrar_service_v1.RegisterMachineRequest{
		Address: h.config.RequestDtos.MachineAddress,
	}

	response, err := client.RegisterMachine(h.t.Context(), &req)

	require.NoError(h.t, err, "Failed to register machine")

	require.NotEmpty(h.t, response.GetMachineId(), "Machine ID should not be empty")
	require.Equal(h.t, h.config.RequestDtos.MachineAddress, response.GetAddress(), "Machine address should match the request")

	return response
}
