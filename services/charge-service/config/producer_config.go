package config

import (
	"fmt"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ChargeProducerConfig struct {
	KafkaConfig *kafka.ConfigMap
	Topic       string
}

func NewChargeProducerConfig() (*ChargeProducerConfig, error) {
	bootstrapServers := os.Getenv("CHARGE_KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		return nil, fmt.Errorf("CHARGE_KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	topic := os.Getenv("CHARGE_KAFKA_TOPIC")
	if topic == "" {
		return nil, fmt.Errorf("CHARGE_KAFKA_TOPIC environment variable is not set")
	}

	kafkaConfig := &kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"client.id":         "charge-service-producer",
		"acks":              "all",
		"retries":           10,
		"batch.size":        16384,
		"linger.ms":         5,
	}

	return &ChargeProducerConfig{
		KafkaConfig: kafkaConfig,
		Topic:       topic,
	}, nil
}
