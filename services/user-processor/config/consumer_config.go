package config

import (
	"fmt"
	"os"
	"time"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
)

func NewUserConsumerConfig() (*kafka.ConsumerConfig, error) {
	bootstrapServers := os.Getenv("USER_KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		return nil, fmt.Errorf("USER_KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	topic := os.Getenv("USER_KAFKA_TOPIC")
	if topic == "" {
		return nil, fmt.Errorf("USER_KAFKA_TOPIC environment variable is not set")
	}

	groupID := os.Getenv("USER_KAFKA_GROUP_ID")
	if groupID == "" {
		return nil, fmt.Errorf("USER_KAFKA_GROUP_ID environment variable is not set")
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
