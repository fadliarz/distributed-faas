package main

import (
	"testing"

	user_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/user-service/v1"
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

func (h *ArrangeHelper) CreateUser(endpoint, password string) (*user_service_v1.CreateUserResponse, error) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := user_service_v1.NewUserServiceClient(conn)

	req := user_service_v1.CreateUserRequest{
		Password: password,
	}

	response, err := client.CreateUser(h.t.Context(), &req)

	return response, err
}
