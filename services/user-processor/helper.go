package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

// Config

type Config struct {
	ShutdownTimeout time.Duration
}

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
	// Application Service
	applicationService := application.NewUserEventHandler(
		application.NewUserProcessorDataMapper(),
		domain.NewUserProcessorDomainService(),
		application.NewUserEventHandlerRepositoryManager(repositoryManager.Cron),
	)

	// Config
	kafkaConfig, err := config.NewUserConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	// Setup Kafka consumer
	deserializer := messaging.NewUserMessageDeserializer()
	processor := messaging.NewUserMessageProcessor(applicationService)

	confluentConsumer, err := kafka.NewConfluentKafkaConsumer(ctx, kafkaConfig, deserializer, processor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	consumer := messaging.NewUserConsumer(confluentConsumer)

	log.Info().Msg("User consumer successfully initialized")

	return consumer, nil
}

// Consumer

func startConsumer(consumer application.UserConsumer, shutdown <-chan os.Signal, timeout time.Duration) {
	// Start consumer in a goroutine
	consumerErr := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting user consumer...")

		consumer.PollAndProcessMessages()
	}()

	// Wait for shutdown signal
	select {
	case err := <-consumerErr:
		log.Fatal().Msgf("consumer failed: %v", err)
	case sig := <-shutdown:
		log.Info().Msgf("Received signal: %s, shutting down...", sig)

		gracefulShutdown(consumer, timeout)
	}
}

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}

func gracefulShutdown(consumer application.UserConsumer, timeout time.Duration) {
	log.Info().Msg("Gracefully stopping consumer...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	done := make(chan struct{})

	go func() {
		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Warn().Msg("Shutdown timeout exceeded")
	case <-done:
		log.Info().Msg("Consumer stopped gracefully")
	}
}
