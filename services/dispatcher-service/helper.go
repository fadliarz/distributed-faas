package main

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/config"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/infrastructure/messaging"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/infrastructure/repository"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dependencies struct {
	invocationConsumer application.InvocationConsumer
}

type RepositoryManager struct {
	Machine application.MachineRepository
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found or failed to load .env file")
	}
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	invocationConsumer, err := setupInvocationConsumer(ctx, repositoryManager)
	if err != nil {
		return nil, fmt.Errorf("failed to setup invocation consumer: %w", err)
	}

	return &Dependencies{
		invocationConsumer: invocationConsumer,
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	machineRepository, err := setupMachineRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup machine repository: %w", err)
	}

	return &RepositoryManager{
		Machine: machineRepository,
	}, nil
}

func setupMachineRepository(ctx context.Context) (application.MachineRepository, error) {
	// Config
	config, err := config.SetupMachineMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load machine mongo config: %w", err)
	}

	// Collection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to MongoDB")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to ping MongoDB")
	}

	collection := client.Database(config.Database).Collection(config.Collection)

	return repository.NewMachineRepository(repository.NewMachineDataAccessMapper(), repository.NewMachineMongoRepository(collection)), nil
}

func setupInvocationConsumer(ctx context.Context, repositoryManager *RepositoryManager) (application.InvocationConsumer, error) {
	// Config
	config, err := config.SetupInvocationConsumerConfig()
	if err != nil {
		return nil, err
	}

	// Confluent Consumer
	handler := application.NewInvocationEventHandler(application.NewInvocationEventHandlerRepositoryManager(repositoryManager.Machine))
	processor := messaging.NewInvocationMessageProcessor(handler)

	confluentConsumer, err := kafka.NewConfluentKafkaConsumer(ctx, config, messaging.NewInvocationMessageDeserializer(), processor)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create ConfluentKafkaConsumer")
	}

	// Consumer
	consumer := messaging.NewInvocationConsumer(confluentConsumer)

	return consumer, nil
}
