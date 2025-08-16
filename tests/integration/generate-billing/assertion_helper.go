package main

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
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

func (ah *AssertionHelper) AssertCronUpdatedInMongoDB(ctx context.Context, client *mongo.Client, userID primitive.ObjectID, expectedUpdatedLastBilled int64) {
	collection := client.Database(ah.config.MongoConfig.CronDatabase).Collection(ah.config.MongoConfig.CronCollection)

	deadline := time.Now().Add(60 * time.Second)

	for time.Now().Before(deadline) {
		var cronEntity CronEntity
		err := collection.FindOne(ctx, bson.M{
			"_id": userID,
		}).Decode(&cronEntity)

		if err != nil || cronEntity.LastBilled != expectedUpdatedLastBilled {
			time.Sleep(2 * time.Second)

			continue
		}

		require.Equal(ah.t, expectedUpdatedLastBilled, cronEntity.LastBilled, "Last billed timestamp should match")

		return
	}

	require.Fail(ah.t, "Failed to find cron entity in MongoDB for user")
}

func (ah *AssertionHelper) AssertUpdatedCronPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, userID primitive.ObjectID, expectedUpdatedLastBilled int64) {
	deadline := time.Now().Add(120 * time.Second)

	for time.Now().Before(deadline) {
		msg, err := consumer.ReadMessage(3 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().Msg("[AssertUpdatedCronPersistedInMongoDB] No message received from Kafka within the timeout period, retrying...")

				continue
			}

			require.NoError(ah.t, err, "Failed to read message from Kafka")
		}

		if msg != nil && len(msg.Value) > 0 {
			var event CronEvent
			err := json.Unmarshal(msg.Value, &event)

			if err != nil || event.LastBilled != expectedUpdatedLastBilled {
				log.Debug().Err(err).Msg("[AssertUpdatedCronPersistedInMongoDB] Failed to unmarshal Kafka message")

				continue
			}

			require.Equal(ah.t, event.LastBilled, expectedUpdatedLastBilled, "Last billed timestamp should match")

			return
		}
	}

	require.Fail(ah.t, "No message received from Kafka within the timeout period")
}
