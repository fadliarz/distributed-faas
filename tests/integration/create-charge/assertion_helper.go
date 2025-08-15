package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/application-service"
	user_service_v1 "github.com/fadliarz/distributed-faas/tests/integration/gen/go/user-service/v1"
	"github.com/rs/zerolog/log"
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

func (ah *AssertionHelper) AssertChargesPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, userID1, userID2 string, serviceID1, serviceID2 string, expectedAggregatedAmount int64) {
	chargeMap := make(map[string]int64)

	deadline := time.Now().Add(120 * time.Second)

	for time.Now().Before(deadline) {
		if len(chargeMap) == 4 {
			return
		}

		msg, err := consumer.ReadMessage(3 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().Msg("[AssertUserPersistedInKafka] No message received from Kafka within the timeout period, retrying...")

				continue
			}

			require.NoError(ah.t, err, "Failed to read message from Kafka")
		}

		if msg != nil && len(msg.Value) > 0 {
			var event ChargeEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Debug().Err(err).Msg("[AssertUserPersistedInKafka] Failed to unmarshal Kafka message")

				continue
			}

			if (event.UserID == userID1 || event.UserID == userID2) && (event.ServiceID == serviceID1 || event.ServiceID == serviceID2) {
				key := fmt.Sprintf("%s:%s", event.UserID, event.ServiceID)

				if _, exist := chargeMap[key]; !exist {
					require.Equal(ah.t, expectedAggregatedAmount, event.AggregatedAmount, "Aggregated amount does not match")

					chargeMap[key] = 1
				}
			}
		}
	}

	require.Fail(ah.t, "No message received from Kafka within the timeout period")
}

func (ah *AssertionHelper) AssertChargePersistedInMongoDB(ctx context.Context, client *mongo.Client, userID, serviceID string, expectedAccumulatedAmount int64) {
	collection := client.Database(ah.config.MongoConfig.ChargeDatabase).Collection(ah.config.MongoConfig.ChargeCollection)

	timestamp := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC).Unix()

	for i := 0; i < 5; i++ {
		var chargeEntity ChargeEntity
		err := collection.FindOne(ctx, bson.M{
			"user_id":    userID,
			"service_id": serviceID,
			"timestamp":  timestamp,
		}).Decode(&chargeEntity)

		if err != nil {
			time.Sleep(2 * time.Second)

			continue
		}

		require.Equal(ah.t, expectedAccumulatedAmount, chargeEntity.AccumulatedAmount, "Accumulated amount should match")

		return
	}

	require.Fail(ah.t, "Failed to find charge entity in MongoDB for user")
}

func (ah *AssertionHelper) AssertUserPersistedInKafka(ctx context.Context, consumer *kafka.Consumer, createUserResponse *user_service_v1.CreateUserResponse) {
	deadline := time.Now().Add(120 * time.Second)

	for time.Now().Before(deadline) {
		msg, err := consumer.ReadMessage(3 * time.Second)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().Msg("[AssertUserPersistedInKafka] No message received from Kafka within the timeout period, retrying...")

				continue
			}

			require.NoError(ah.t, err, "Failed to read message from Kafka")
		}

		if msg != nil && len(msg.Value) > 0 {
			var event application.UserEvent
			err := json.Unmarshal(msg.Value, &event)
			if err != nil {
				log.Debug().Err(err).Msg("[AssertUserPersistedInKafka] Failed to unmarshal Kafka message")

				continue
			}

			require.Equal(ah.t, createUserResponse.UserId, event.UserID, "User ID does not match")

			return
		}
	}

	require.Fail(ah.t, "No message received from Kafka within the timeout period")
}
