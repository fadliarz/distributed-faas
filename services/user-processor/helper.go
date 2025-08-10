package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/user-processor/config"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
	"github.com/fadliarz/distributed-faas/services/user-processor/infrastructure/messaging"
	"github.com/fadliarz/distributed-faas/services/user-processor/infrastructure/repository"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Env

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using default configuration")
	}
}

// Dependencies

type Dependencies struct {
	consumer application.UserConsumer
}

type RepositoryManager struct {
	Cron application.CronRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	userConsumer, err := setupUserConsumer(ctx, repositoryManager)
	if err != nil {
		return nil, fmt.Errorf("failed to setup user consumer: %w", err)
	}

	return &Dependencies{
		consumer: userConsumer,
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	cronRepository, err := setupCronRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup cron repository: %w", err)
	}

	return &RepositoryManager{
		Cron: cronRepository,
	}, nil
}

func setupCronRepository(ctx context.Context) (application.CronRepository, error) {
	// Config
	config, err := config.NewCronMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load cron mongo config: %w", err)
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

	return repository.NewCronRepository(repository.NewCronDataAccessMapper(), repository.NewCronMongoRepository(collection)), nil
}

func setupUserConsumer(ctx context.Context, repositoryManager *RepositoryManager) (application.UserConsumer, error) {
	// Config
	config, err := config.NewUserConsumerConfig()
	if err != nil {
		return nil, err
	}

	// Confluent Consumer
	handler := application.NewUserEventHandler(application.NewUserProcessorDataMapper(), domain.NewUserProcessorDomainService(), application.NewUserEventHandlerRepositoryManager(repositoryManager.Cron))
	processor := messaging.NewUserMessageProcessor(handler)

	confluentConsumer, err := kafka.NewConfluentKafkaConsumer(ctx, config, messaging.NewUserMessageDeserializer(), processor)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create ConfluentKafkaConsumer")
	}

	// Consumer
	consumer := messaging.NewUserConsumer(confluentConsumer)

	return consumer, nil
}

// Servers

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}
