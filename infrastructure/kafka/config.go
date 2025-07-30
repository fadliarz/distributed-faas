package kafka

import "time"

type ConsumerConfig struct {
	Basic      *ConsumerBasicConfig
	Processing *ConsumerProcessingConfig
}

type ConsumerBasicConfig struct {
	BootstrapServers string
	Topic            string
	GroupID          string
	PollTimeout      time.Duration
}

type ConsumerProcessingConfig struct {
	NumWorkers int
}
