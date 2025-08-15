package main

import (
	"fmt"
	"testing"

	charge_service_v1 "github.com/fadliarz/distributed-faas/services/charge-service/gen/go/charge-service/v1"
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

func (h *ArrangeHelper) CreateCharges(endpoint string, requests []*charge_service_v1.CreateChargeRequest) error {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(h.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := charge_service_v1.NewChargeServiceClient(conn)

	for _, request := range requests {
		_, err := client.CreateCharge(h.t.Context(), request)
		if err != nil {
			return fmt.Errorf("failed to create charge: %w", err)
		}
	}

	return nil
}
