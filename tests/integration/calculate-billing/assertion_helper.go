package main

import (
	"context"
	"testing"
	"time"

	billing_service_v1 "github.com/fadliarz/distributed-faas/services/billing-service/gen/go/v1"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AssertionHelper struct {
	t      *testing.T
	config *TestConfig
}

func NewAssertionHelper(t *testing.T, config *TestConfig) *AssertionHelper {
	return &AssertionHelper{
		t:      t,
		config: config,
	}
}

func (ah *AssertionHelper) AssertBillingPersistedInMongoDB(ctx context.Context, client *mongo.Client, userID string, expectedAmount int64) {
	collection := client.Database(ah.config.MongoConfig.BillingDatabase).Collection(ah.config.MongoConfig.BillingCollection)

	for i := 0; i < 5; i++ {
		var billingEntity BillingEntity
		err := collection.FindOne(ctx, bson.M{
			"user_id": userID,
		}).Decode(&billingEntity)

		if err != nil {
			time.Sleep(2 * time.Second)

			continue
		}

		require.Equal(ah.t, expectedAmount, billingEntity.Amount, "Amount should match")

		return
	}

	require.Fail(ah.t, "Failed to find billing entity in MongoDB for user")
}

func (ah *AssertionHelper) AssertGetBillingSuccessfully(ctx context.Context, endpoint string, userID string, expectedLastBilled, expectedAmount int64) {
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(ah.t, err, "Failed to connect to gRPC server")

	defer conn.Close()

	client := billing_service_v1.NewBillingServiceClient(conn)

	billing, err := client.GetBilling(ctx, &billing_service_v1.GetBillingRequest{
		UserId: userID,
	})
	require.NoError(ah.t, err, "Failed to call GetBilling")

	require.Equal(ah.t, expectedLastBilled, billing.LastBilled, "LastBilled should match")
	require.Equal(ah.t, expectedAmount, billing.Amount, "Amount should match")
}
