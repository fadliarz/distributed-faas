package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ChargeKafkaConfig struct {
	BootstrapServers string
	GroupID          string
	Topic            string
	PollTimeout      time.Duration
	MaxBatchSize     int
}

func NewChargeKafkaConfig() (*ChargeKafkaConfig, error) {
	bootstrapServers := os.Getenv("CHARGE_KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		return nil, errors.New("CHARGE_KAFKA_BOOTSTRAP_SERVERS environment variable is not set")
	}

	topic := os.Getenv("CHARGE_KAFKA_TOPIC")
	if topic == "" {
		return nil, errors.New("CHARGE_KAFKA_TOPIC environment variable is not set")
	}

	groupID := os.Getenv("CHARGE_KAFKA_GROUP_ID")
	if groupID == "" {
		return nil, errors.New("CHARGE_KAFKA_GROUP_ID environment variable is not set")
	}

	// Default poll duration to 60 seconds as specified in requirements
	pollDurationStr := os.Getenv("CHARGE_KAFKA_POLL_DURATION_SEC")
	pollDuration := 60 * time.Second
	if pollDurationStr != "" {
		if seconds, err := strconv.Atoi(pollDurationStr); err == nil {
			pollDuration = time.Duration(seconds) * time.Second
		}
	}

	// Default batch size to 1000 as specified in requirements
	maxBatchSizeStr := os.Getenv("CHARGE_KAFKA_MAX_BATCH_SIZE")
	maxBatchSize := 1000
	if maxBatchSizeStr != "" {
		if size, err := strconv.Atoi(maxBatchSizeStr); err == nil && size > 0 {
			maxBatchSize = size
		}
	}

	return &ChargeKafkaConfig{
		BootstrapServers: bootstrapServers,
		GroupID:          groupID,
		Topic:            topic,
		PollTimeout:      pollDuration,
		MaxBatchSize:     maxBatchSize,
	}, nil
}

func (c *ChargeKafkaConfig) ToKafkaConfigMap() kafka.ConfigMap {
	return kafka.ConfigMap{
		"bootstrap.servers":               c.BootstrapServers,
		"group.id":                        c.GroupID,
		"auto.offset.reset":               "earliest",
		"enable.auto.commit":              false, // Disable auto commit for manual batch commit
		"session.timeout.ms":              30000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
	}
}
