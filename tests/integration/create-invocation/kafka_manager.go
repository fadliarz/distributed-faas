package main

import (
	"context"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/rs/zerolog/log"
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
	if err := km.createTopicsIfNotExist(bootstrapServer, []string{km.config.KafkaConfig.InvocationTopic, km.config.KafkaConfig.CheckpointTopic}); err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	if err := km.setupInvocationConsumer(bootstrapServer); err != nil {
		return fmt.Errorf("failed to setup invocation consumers: %w", err)
	}

	if err := km.setupCheckpointConsumer(bootstrapServer); err != nil {
		return fmt.Errorf("failed to setup checkpoint consumers: %w", err)
	}

	return nil
}

func (km *KafkaManager) createTopicsIfNotExist(bootstrapServer string, topics []string) error {
	adminClient, err := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
	})
	if err != nil {
		return fmt.Errorf("failed to create admin client: %w", err)
	}
	defer adminClient.Close()

	metadata, err := adminClient.GetMetadata(nil, false, 5000)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get metadata, proceeding to create all topics")
	}

	var topicsToCreate []kafka.TopicSpecification
	for _, topic := range topics {
		if metadata != nil {
			if _, exists := metadata.Topics[topic]; exists {
				log.Debug().Msgf("Topic %s already exists", topic)

				continue
			}
		}

		topicsToCreate = append(topicsToCreate, kafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
	}

	if len(topicsToCreate) == 0 {
		log.Debug().Msg("All topics already exist")

		return nil
	}

	results, err := adminClient.CreateTopics(km.ctx, topicsToCreate)
	if err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	for _, result := range results {
		if result.Error.Code() != kafka.ErrNoError {
			if result.Error.Code() == kafka.ErrTopicAlreadyExists {
				continue
			}

			return fmt.Errorf("failed to create topic %s: %v", result.Topic, result.Error)
		}
	}

	return nil
}

func (km *KafkaManager) setupInvocationConsumer(bootstrapServer string) error {
	var err error

	km.Consumers.Invocation, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServer,
		"group.id":           "group-1-invocation",
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

func (km *KafkaManager) setupCheckpointConsumer(bootstrapServer string) error {
	var err error

	km.Consumers.Checkpoint, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  bootstrapServer,
		"group.id":           "group-1-checkpoint",
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
