package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fadliarz/distributed-faas/services/charge-service/config"
	"github.com/fadliarz/distributed-faas/services/charge-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/charge-service/domain/domain-core"
	"github.com/rs/zerolog/log"
)

type ChargeEvent struct {
	UserID           string `json:"user_id"`
	ServiceID        string `json:"service_id"`
	AggregatedAmount int64  `json:"aggregated_amount"`
}

type ChargeProducerImpl struct {
	producer *kafka.Producer
	config   *config.ChargeProducerConfig
}

func NewChargeProducer(ctx context.Context, config *config.ChargeProducerConfig) (application.ChargeProducer, error) {
	producer, err := kafka.NewProducer(config.KafkaConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	chargeProducer := &ChargeProducerImpl{
		producer: producer,
		config:   config,
	}

	go chargeProducer.handleDeliveryReports(ctx)

	return chargeProducer, nil
}

func (p *ChargeProducerImpl) SendAggregatedCharges(ctx context.Context, aggregates []*domain.ChargeAggregate) error {
	for _, aggregate := range aggregates {
		message := &ChargeEvent{
			UserID:           aggregate.UserID.String(),
			ServiceID:        aggregate.ServiceID.String(),
			AggregatedAmount: aggregate.AggregatedAmount.Int64(),
		}

		// Convert to JSON
		jsonData, err := json.Marshal(message)
		if err != nil {
			return fmt.Errorf("failed to marshal charge aggregate to JSON: %w", err)
		}

		key := aggregate.GetKey()

		// Send to Kafka
		err = p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     &p.config.Topic,
				Partition: kafka.PartitionAny,
			},
			Key:   []byte(key),
			Value: jsonData,
		}, nil)

		if err != nil {
			return fmt.Errorf("failed to send charge aggregate with key %s: %w", key, err)
		}
	}

	return nil
}

func (p *ChargeProducerImpl) handleDeliveryReports(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("Stopping Kafka delivery report handler")
			return
		case e := <-p.producer.Events():
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error().
						Err(ev.TopicPartition.Error).
						Str("topic", *ev.TopicPartition.Topic).
						Int32("partition", ev.TopicPartition.Partition).
						Msg("Failed to deliver message to Kafka")
				} else {
					log.Debug().
						Str("topic", *ev.TopicPartition.Topic).
						Int32("partition", ev.TopicPartition.Partition).
						Int64("offset", int64(ev.TopicPartition.Offset)).
						Msg("Message delivered to Kafka")
				}
			case kafka.Error:
				log.Error().
					Err(ev).
					Msg("Kafka producer error")
			}
		}
	}
}

func (p *ChargeProducerImpl) Close() {
	if p.producer != nil {
		p.producer.Close()
	}
}
