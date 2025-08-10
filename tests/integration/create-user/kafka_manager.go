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
	User *kafka.Consumer
	Cron *kafka.Consumer
}

func NewKafkaManager(ctx context.Context, config *TestConfig) *KafkaManager {
	return &KafkaManager{
		ctx:       ctx,
		config:    config,
		Consumers: &Consumers{},
	}
}

func (km *KafkaManager) SetupConsumers(bootstrapServer string) error {
	if err := km.createTopicsIfNotExist(bootstrapServer, []string{km.config.KafkaConfig.Topics.User, km.config.KafkaConfig.Topics.Cron}); err != nil {
		return fmt.Errorf("failed to create topics: %w", err)
	}

	if err := km.setupUserConsumer(bootstrapServer); err != nil {
		return fmt.Errorf("failed to setup user consumer: %w", err)
	}

	if err := km.setupCronConsumer(bootstrapServer); err != nil {
		return fmt.Errorf("failed to setup cron consumer: %w", err)
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

func (km *KafkaManager) setupUserConsumer(bootstrapServer string) error {
	var err error

	km.Consumers.User, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          km.config.KafkaConfig.ConsumerGroups.User,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return fmt.Errorf("failed to create user consumer: %w", err)
	}

	err = km.Consumers.User.Subscribe(km.config.KafkaConfig.Topics.User, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to user topic: %w", err)
	}

	log.Info().Msgf("User consumer subscribed to topic: %s", km.config.KafkaConfig.Topics.User)

	return nil
}

func (km *KafkaManager) setupCronConsumer(bootstrapServer string) error {
	var err error

	km.Consumers.Cron, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          km.config.KafkaConfig.ConsumerGroups.Cron,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return fmt.Errorf("failed to create cron consumer: %w", err)
	}

	err = km.Consumers.Cron.Subscribe(km.config.KafkaConfig.Topics.Cron, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to cron topic: %w", err)
	}

	log.Info().Msgf("Cron consumer subscribed to topic: %s", km.config.KafkaConfig.Topics.Cron)

	return nil
}
