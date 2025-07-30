package config

import (
	"fmt"
	"os"
	"time"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
)

func SetupInvocationConsumerConfig() (*kafka.ConsumerConfig, error) {
	bootstrapServers := os.Getenv("INVOCATION_KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		return nil, fmt.Errorf("INVOCATION_KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	topic := os.Getenv("INVOCATION_KAFKA_TOPIC")
	if topic == "" {
		return nil, fmt.Errorf("INVOCATION_KAFKA_TOPIC environment variable is not set")
	}

	groupID := os.Getenv("INVOCATION_KAFKA_GROUP_ID")
	if groupID == "" {
		return nil, fmt.Errorf("INVOCATION_KAFKA_GROUP_ID environment variable is not set")
	}

	ConsumerBasicConfig := &kafka.ConsumerBasicConfig{
		BootstrapServers: bootstrapServers,
		Topic:            topic,
		GroupID:          groupID,
		PollTimeout:      3 * time.Second,
	}

	ConsumerProcessingConfig := &kafka.ConsumerProcessingConfig{
		NumWorkers: 3,
	}

	return &kafka.ConsumerConfig{
		Basic:      ConsumerBasicConfig,
		Processing: ConsumerProcessingConfig,
	}, nil
}
