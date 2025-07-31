package main

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaManager struct {
	ctx    context.Context
	config *TestConfig

	Consumers *Consumers
}

type Consumers struct {
	Invocation *kafka.Consumer
	Checkpoint *kafka.Consumer
}

func NewKafkaManager(ctx context.Context, config *TestConfig) *KafkaManager {
	return &KafkaManager{
		ctx:       ctx,
		config:    config,
		Consumers: &Consumers{},
	}
}

func (km *KafkaManager) SetupConsumers(bootstrapServer string) error {

	if err := km.setupInvocationConsumers(bootstrapServer); err != nil {
		return fmt.Errorf("failed to setup invocation consumers: %w", err)
	}

	if err := km.setupCheckpointConsumers(bootstrapServer); err != nil {
		return fmt.Errorf("failed to setup checkpoint consumers: %w", err)
	}

	return nil

}

func (km *KafkaManager) setupInvocationConsumers(bootstrapServer string) error {
	var err error

	km.Consumers.Invocation, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServer,
		"group.id":           "group-1",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": km.config.KafkaConfig.AutoCommit,
	})
	if err != nil {
		return fmt.Errorf("failed to create a Kafka consumer: %w", err)
	}

	err = km.Consumers.Invocation.SubscribeTopics([]string{km.config.KafkaConfig.InvocationTopic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to Kafka topics: %w", err)
	}

	return nil
}

func (km *KafkaManager) setupCheckpointConsumers(bootstrapServer string) error {
	var err error

	km.Consumers.Checkpoint, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServer,
		"group.id":           "group-1",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": km.config.KafkaConfig.AutoCommit,
	})
	if err != nil {
		return fmt.Errorf("failed to create a Kafka consumer: %w", err)
	}

	err = km.Consumers.Checkpoint.SubscribeTopics([]string{km.config.KafkaConfig.CheckpointTopic}, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to Kafka topics: %w", err)
	}

	return nil
}
