package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/config"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/application-service"
	"github.com/rs/zerolog/log"
)

// Custom Kafka Consumer Implementation

type ChargeConsumerImpl struct {
	consumer     *kafka.Consumer
	config       *config.AccumulatorKafkaConfig
	eventHandler *application.ChargeEventHandler
}

func NewChargeConsumer(config *config.AccumulatorKafkaConfig, eventHandler *application.ChargeEventHandler) (application.ChargeConsumer, error) {
	configMap := config.ToKafkaConfigMap()

	consumer, err := kafka.NewConsumer(&configMap)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	err = consumer.Subscribe(config.Topic, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", config.Topic, err)
	}

	return &ChargeConsumerImpl{
		consumer:     consumer,
		config:       config,
		eventHandler: eventHandler,
	}, nil
}

func (c *ChargeConsumerImpl) PollAndProcessMessages() {
	log.Debug().
		Str("topic", c.config.Topic).
		Str("group_id", c.config.GroupID).
		Dur("poll_timeout", c.config.PollTimeout).
		Int("max_batch_size", c.config.MaxBatchSize).
		Msg("Starting Kafka consumer polling")

	defer c.consumer.Close()

	for {
		c.pollAndProcessBatch()
	}
}

func (c *ChargeConsumerImpl) pollAndProcessBatch() {
	ctx := context.Background()

	var message *kafka.Message
	var events []*application.ChargeEvent

	deadline := time.Now().Add(60 * time.Second)

	for len(events) < c.config.MaxBatchSize && time.Now().Before(deadline) {
		msg, err := c.consumer.ReadMessage(c.config.PollTimeout)
		if err != nil {
			if err.(kafka.Error).Code() == kafka.ErrTimedOut {
				log.Debug().
					Int("batch_size", len(events)).
					Msg("Poll timeout reached, processing collected messages")

				continue
			}

			log.Error().Err(err).Msg("Error polling messages")

			continue
		}

		// Deserialize the message
		event, err := c.deserializeMessage(msg.Value)
		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to deserialize message, skipping")

			continue
		}

		events = append(events, event)
		message = msg
	}

	for message != nil {
		_, err := c.consumer.CommitMessage(message)
		if err != nil {
			log.Error().
				Err(err).
				Int("batch_size", len(events)).
				Msg("Failed to commit batch, retrying...")

			time.Sleep(5 * time.Second)

			continue
		}

		break
	}

	if len(events) > 0 {
		err := c.processBatch(ctx, events)
		if err != nil {
			log.Error().
				Err(err).
				Int("batch_size", len(events)).
				Msg("Failed to process message batch")

			return // Continue processing despite error (requirement: it's okay to lose messages)
		}
	}
}

func (c *ChargeConsumerImpl) processBatch(ctx context.Context, events []*application.ChargeEvent) error {
	return c.eventHandler.HandleChargeEventBatch(ctx, events)
}

func (c *ChargeConsumerImpl) deserializeMessage(data []byte) (*application.ChargeEvent, error) {
	var event application.ChargeEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to deserialize message: %w", err)
	}

	return &event, nil
}
