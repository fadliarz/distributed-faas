package main

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"go.mongodb.org/mongo-driver/mongo"
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

func (h *ArrangeHelper) CreateCharges(ctx context.Context, client *mongo.Client, entities []*ChargeEntity) error {
	collection := client.Database(h.config.MongoConfig.ChargeDatabase).Collection(h.config.MongoConfig.ChargeCollection)

	for _, entity := range entities {
		_, err := collection.InsertOne(ctx, entity)
		if err != nil {
			return fmt.Errorf("failed to create charge: %w", err)
		}
	}

	return nil
}

func (h *ArrangeHelper) CreateCronEvents(ctx context.Context, producer *kafka.Producer, events []*CronEvent) error {
	for _, event := range events {
		jsonData, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal cron event: %w", err)
		}

		err = producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &h.config.KafkaConfig.Topics.Cron,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(fmt.Sprintf("%s:%d", event.UserID, event.LastBilled)),
			Value: jsonData,
		}, nil)
		if err != nil {
			return fmt.Errorf("failed to produce cron event: %w", err)
		}
	}

	return nil
}
