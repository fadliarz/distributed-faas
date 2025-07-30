package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/config"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/infrastructure/messaging"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/infrastructure/repository"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config

type Config struct {
	ShutdownTimeout time.Duration
}

// Env

func loadEnv(config *Config) {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using default configuration")
	}
}

// Dependencies

type Dependencies struct {
	consumer application.CheckpointConsumer
}

type RepositoryManager struct {
	Invocation application.InvocationRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	checkpointConsumer, err := setupCheckpointConsumer(ctx, repositoryManager)
	if err != nil {
		return nil, fmt.Errorf("failed to setup checkpoint consumer: %w", err)
	}

	return &Dependencies{
		consumer: checkpointConsumer,
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	invocationRepository, err := setupInvocationRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup invocation repository: %w", err)
	}

	return &RepositoryManager{
		Invocation: invocationRepository,
	}, nil
}

func setupInvocationRepository(ctx context.Context) (application.InvocationRepository, error) {
	// Config
	config, err := config.NewInvocationMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load invocation mongo config: %w", err)
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

	log.Info().Msgf("Ping MongoDB successful for database %s and collection %s", config.Database, config.Collection)

	collection := client.Database(config.Database).Collection(config.Collection)

	return repository.NewInvocationRepository(repository.NewInvocationMongoRepository(collection)), nil
}

func setupCheckpointConsumer(ctx context.Context, repositoryManager *RepositoryManager) (application.CheckpointConsumer, error) {
	// Config
	config, err := config.NewCheckpointConsumerConfig()
	if err != nil {
		return nil, err
	}

	// Confluent Consumer
	handler := application.NewCheckpointEventHandler(application.NewCheckpointProcessorDataMapper(), application.NewCheckpointEventHandlerRepositoryManager(repositoryManager.Invocation))
	processor := messaging.NewCheckpointMessageProcessor(handler)

	confluentConsumer, err := kafka.NewConfluentKafkaConsumer(ctx, config, messaging.NewCheckpointMessageDeserializer(), processor)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create ConfluentKafkaConsumer")
	}

	// Consumer
	consumer := messaging.NewCheckpointConsumer(confluentConsumer)

	return consumer, nil
}

// Servers

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}
