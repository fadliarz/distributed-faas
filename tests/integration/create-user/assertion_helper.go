package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	up_repository "github.com/fadliarz/distributed-faas/services/user-processor/infrastructure/repository"
	us_repository "github.com/fadliarz/distributed-faas/services/user-service/infrastructure/repository"
	user_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/user-service/v1"
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

func (ah *AssertionHelper) AssertUserCreatedSuccessfully(createUserResponse *user_service_v1.CreateUserResponse, err error) {
	require.NoError(ah.t, err, "Failed to create user")
	require.NotEmpty(ah.t, createUserResponse.GetUserId(), "User ID should not be empty")
}

func (ah *AssertionHelper) AssertUserPersistedInMongoDB(ctx context.Context, client *mongo.Client, userID string, expectedPassword string) {
	userObjectID, err := primitive.ObjectIDFromHex(userID)
	require.NoError(ah.t, err, "Failed to convert user ID to ObjectID")

	collection := client.Database(ah.config.MongoConfig.UserDatabase).Collection(ah.config.MongoConfig.UserCollection)

	var userEntity us_repository.UserEntity
	err = collection.FindOne(ctx, bson.M{"_id": userObjectID}).Decode(&userEntity)
	require.NoError(ah.t, err, "Failed to find user in MongoDB")

	require.Equal(ah.t, userObjectID, userEntity.UserID, "User ID should match")
	require.Equal(ah.t, expectedPassword, userEntity.Password, "Password should match")
}

func (ah *AssertionHelper) AssertUserPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, createUserResponse *user_service_v1.CreateUserResponse) {
	timeout := time.NewTimer(ah.config.KafkaConfig.Timeout)
	defer timeout.Stop()

	for {
		select {
		case <-timeout.C:
			ah.t.Fatal("Timeout waiting for user event in Kafka")
		default:
			message, err := consumer.ReadMessage(1 * time.Second)
			if err != nil {
				continue
			}

			var userEvent UserEvent
			err = json.Unmarshal(message.Value, &userEvent)
			require.NoError(ah.t, err, "Failed to unmarshal user event")

			if userEvent.After != nil && userEvent.After.UserID == createUserResponse.UserId {
				log.Info().Msgf("User event found in Kafka: %s", createUserResponse.UserId)
				return
			}
		}
	}
}

func (ah *AssertionHelper) AssertCronJobPersistedInMongoDB(ctx context.Context, client *mongo.Client, userID string) {
	collection := client.Database(ah.config.MongoConfig.CronDatabase).Collection(ah.config.MongoConfig.CronCollection)

	timeout := time.NewTimer(ah.config.KafkaConfig.Timeout)
	defer timeout.Stop()

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout.C:
			ah.t.Fatal("Timeout waiting for cron job to be persisted in MongoDB")
		case <-ticker.C:
			var cronEntity up_repository.CronEntity
			err := collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&cronEntity)
			if err == nil {
				require.Equal(ah.t, userID, cronEntity.UserID, "User ID should match in cron job")
				log.Info().Msgf("Cron job found in MongoDB for user: %s", userID)
				return
			}
		}
	}
}
