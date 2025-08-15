package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
