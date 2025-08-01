package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/application-service"
	fs_repository "github.com/fadliarz/distributed-faas/services/function-service/infrastructure/repository"
	is_repository "github.com/fadliarz/distributed-faas/services/invocation-service/infrastructure/repository"
	invocation_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/invocation-service/v1"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (ah *AssertionHelper) AssertFunctionPersistedInFunctionMongoDB(ctx context.Context, client *mongo.Client, functionID string) {
	collection := client.Database(ah.config.MongoConfig.FunctionDatabase).Collection(ah.config.MongoConfig.FunctionCollection)

	objectID, err := primitive.ObjectIDFromHex(functionID)
	require.NoError(ah.t, err, "Failed to convert Function ID to ObjectID")

	var function fs_repository.FunctionEntity
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&function)
	require.NoError(ah.t, err, "Failed to find function in MongoDB")

	require.NotEmpty(ah.t, function.FunctionID, "Function ID should not be empty")
	require.NotEmpty(ah.t, function.UserID, "User ID should not be empty")
	require.Empty(ah.t, function.SourceCodeURL, "")
}

func (ah *AssertionHelper) AssertInvocationPersistedInMongoDB(ctx context.Context, client *mongo.Client, invocationID string) {
	collection := client.Database(ah.config.MongoConfig.InvocationDatabase).Collection(ah.config.MongoConfig.InvocationCollection)

	objectID, err := primitive.ObjectIDFromHex(invocationID)
	require.NoError(ah.t, err, "Failed to convert Invocation ID to ObjectID")

	var invocation is_repository.InvocationEntity
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invocation)
	require.NoError(ah.t, err, "Failed to find invocation in MongoDB")

	require.NotEmpty(ah.t, invocation.InvocationID, "Invocation ID should not be empty")
	require.NotEmpty(ah.t, invocation.FunctionID, "Function ID should not be empty")
	require.NotEmpty(ah.t, invocation.UserID, "User ID should not be empty")
	require.Empty(ah.t, invocation.SourceCodeURL, "Source code URL should be empty")
	require.Empty(ah.t, invocation.OutputURL, "Output URL should be empty")
	require.False(ah.t, invocation.IsRetry, "Invocation should not be a retry")
	require.Greater(ah.t, invocation.Timestamp, int64(0), "Timestamp should be greater than 0")
}

func (ah *AssertionHelper) AssertInvocationPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, createInvocationResponse *invocation_service_v1.CreateInvocationResponse) {
	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		msg, err := consumer.ReadMessage(3 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().Msg("No message received from Kafka within the timeout period, retrying...")

				continue
			}

			require.NoError(ah.t, err, "Failed to read message from Kafka")
		}

		if msg != nil && len(msg.Value) > 0 {
			var event application.InvocationCreatedEvent
			err := json.Unmarshal(msg.Value, &event)

			require.NoError(ah.t, err, "Failed to unmarshal message from Kafka")
			require.Equal(ah.t, createInvocationResponse.InvocationId, event.InvocationID, "Invocation ID does not match")
			require.Equal(ah.t, createInvocationResponse.FunctionId, event.FunctionID, "Function ID does not match")
			require.Equal(ah.t, createInvocationResponse.SourceCodeUrl, event.SourceCodeURL, "Source code URL does not match")
			require.Less(ah.t, int64(0), event.Timestamp, "Timestamp should be greater than 0")
			require.False(ah.t, event.IsRetry, "Invocation should not be a retry")

			return
		}
	}

	require.Fail(ah.t, "No message received from Kafka within the timeout period")
}
