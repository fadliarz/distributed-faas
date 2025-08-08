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
	m_repository "github.com/fadliarz/distributed-faas/services/machine/infrastructure/repository"
	invocation_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/invocation-service/v1"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (ah *AssertionHelper) AssertInvocationCreatedSuccessfully(createInvocationRespose *invocation_service_v1.CreateInvocationResponse, err error) {
	require.NoError(ah.t, err, "Failed to create invocation")
	require.NotEmpty(ah.t, createInvocationRespose.GetInvocationId(), "Invocation ID should not be empty")
}

func (ah *AssertionHelper) AssertInvocationUnauthorized(createInvocationRespose *invocation_service_v1.CreateInvocationResponse, err error) {
	require.Nil(ah.t, createInvocationRespose, "Create invocation response should be nil")
	require.Error(ah.t, err, "Expected an error when creating invocation")

	st, ok := status.FromError(err)

	require.True(ah.t, ok, "Error should be a gRPC status error")
	require.Equal(ah.t, codes.PermissionDenied.String(), st.Code().String(), "Expected PERMISSION_DENIED status code")
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
	require.NotEmpty(ah.t, invocation.SourceCodeURL, "Source code URL should not be empty")
	require.Empty(ah.t, invocation.OutputURL, "Output URL should be empty")
	require.Greater(ah.t, invocation.Timestamp, int64(0), "Timestamp should be greater than 0")
}

func (ah *AssertionHelper) AssertInvocationPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, createInvocationResponse *invocation_service_v1.CreateInvocationResponse) {
	deadline := time.Now().Add(120 * time.Second)

	for time.Now().Before(deadline) {
		msg, err := consumer.ReadMessage(3 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().Msg("[AssertInvocationPersistedInKafka] No message received from Kafka within the timeout period, retrying...")

				continue
			}

			require.NoError(ah.t, err, "Failed to read message from Kafka")
		}

		if msg != nil && len(msg.Value) > 0 {
			var event application.InvocationCreatedEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				continue
			}

			require.Equal(ah.t, createInvocationResponse.InvocationId, event.InvocationID, "Invocation ID does not match")
			require.Equal(ah.t, createInvocationResponse.FunctionId, event.FunctionID, "Function ID does not match")
			require.Equal(ah.t, createInvocationResponse.SourceCodeUrl, event.SourceCodeURL, "Source code URL does not match")
			require.Less(ah.t, int64(0), event.Timestamp, "Timestamp should be greater than 0")

			return
		}
	}

	require.Fail(ah.t, "No message received from Kafka within the timeout period")
}

func (ah *AssertionHelper) AssertCheckpointPersistedInMongoDB(ctx context.Context, client *mongo.Client, checkpointID string) {
	var err error

	collection := client.Database(ah.config.MongoConfig.CheckpointDatabase).Collection(ah.config.MongoConfig.CheckpointCollection)

	objectID, err := primitive.ObjectIDFromHex(checkpointID)
	require.NoError(ah.t, err, "Failed to convert Checkpoint ID to ObjectID")

	for i := 0; i < 10000; i++ {
		var checkpoint m_repository.CheckpointEntity

		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&checkpoint)
		if err != nil {
			log.Debug().Err(err).Msgf("Failed to find checkpoint in MongoDB, attempt %d", i+1)

			time.Sleep(2 * time.Second)

			continue
		}

		if checkpoint.OutputURL == "" {
			log.Debug().Msgf("Checkpoint Output URL is empty, retrying... (attempt %d)", i+1)

			time.Sleep(2 * time.Second)

			continue
		}

		require.NoError(ah.t, err, "Failed to find checkpoint in MongoDB")

		require.Equal(ah.t, checkpointID, checkpoint.CheckpointID.Hex(), "Checkpoint ID does not match")
		require.Equal(ah.t, checkpoint.Status, "SUCCESS", "Checkpoint status should be SUCCESS")
		require.NotEmpty(ah.t, checkpoint.OutputURL, "Output URL should not be empty")

		return
	}

	require.Fail(ah.t, "Failed to find checkpoint in MongoDB after multiple attempts")
}

func (ah *AssertionHelper) AssertCheckpointPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, createInvocationResponse *invocation_service_v1.CreateInvocationResponse) {
	deadline := time.Now().Add(120 * time.Second)

	for time.Now().Before(deadline) {
		msg, err := consumer.ReadMessage(3 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().Msg("[AssertCheckpointPersistedInKafka] No message received from Kafka within the timeout period, retrying...")

				continue
			}

			require.NoError(ah.t, err, "Failed to read message from Kafka")
		}

		if msg != nil && len(msg.Value) > 0 {
			var event CheckpointEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				continue
			}

			require.Equal(ah.t, createInvocationResponse.InvocationId, event.CheckpointID, "Checkpoint ID does not match")
			require.NotEmpty(ah.t, event.OutputURL, "Output URL should not be empty")

			return
		}
	}

	require.Fail(ah.t, "No message received from Kafka within the timeout period")
}

func (ah *AssertionHelper) AssertInvocationSuccessInMongoDB(ctx context.Context, client *mongo.Client, invocationID string) {
	collection := client.Database(ah.config.MongoConfig.InvocationDatabase).Collection(ah.config.MongoConfig.InvocationCollection)

	objectID, err := primitive.ObjectIDFromHex(invocationID)
	require.NoError(ah.t, err, "Failed to convert Invocation ID to ObjectID")

	for i := 0; i < 100; i++ {
		var invocation is_repository.InvocationEntity

		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invocation)
		require.NoError(ah.t, err, "Failed to find invocation in MongoDB")

		if invocation.Status != "SUCCESS" || invocation.OutputURL == "" {
			log.Debug().Msgf("Invocation status is not SUCCESS or Output URL is empty, retrying... (attempt %d)", i+1)

			time.Sleep(2 * time.Second)

			continue
		}

		require.Equal(ah.t, "SUCCESS", invocation.Status, "Invocation status should be SUCCESS")
		require.NotEmpty(ah.t, invocation.OutputURL, "Output URL should not be empty")

		return
	}
}
