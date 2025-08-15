package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type CronKafkaConfig struct {
	BootstrapServers string
	Topic            string
	GroupID          string
	PollTimeout      time.Duration
	NumWorkers       int
}

func NewCronKafkaConfig() (*CronKafkaConfig, error) {
	bootstrapServers := os.Getenv("CRON_KAFKA_BOOTSTRAP_SERVERS")
	if bootstrapServers == "" {
		return nil, fmt.Errorf("missing CRON_KAFKA_BOOTSTRAP_SERVERS environment variable")
	}

	topic := os.Getenv("CRON_KAFKA_TOPIC")
	if topic == "" {
		return nil, fmt.Errorf("missing CRON_KAFKA_TOPIC environment variable")
	}

	groupID := os.Getenv("CRON_KAFKA_GROUP_ID")
	if groupID == "" {
		return nil, fmt.Errorf("missing CRON_KAFKA_GROUP_ID environment variable")
	}

	// Optional: Poll timeout with default
	pollTimeoutStr := os.Getenv("CRON_KAFKA_POLL_TIMEOUT")
	pollTimeout := 30 * time.Second // default
	if pollTimeoutStr != "" {
		if parsed, err := time.ParseDuration(pollTimeoutStr); err == nil {
			pollTimeout = parsed
		}
	}

	// Optional: Number of workers with default
	numWorkersStr := os.Getenv("CRON_KAFKA_NUM_WORKERS")
	numWorkers := 5 // default
	if numWorkersStr != "" {
		if parsed, err := strconv.Atoi(numWorkersStr); err == nil {
			numWorkers = parsed
		}
	}

	return &CronKafkaConfig{
		BootstrapServers: bootstrapServers,
		Topic:            topic,
		GroupID:          groupID,
		PollTimeout:      pollTimeout,
		NumWorkers:       numWorkers,
	}, nil
}
