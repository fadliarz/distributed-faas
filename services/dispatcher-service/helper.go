package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	invocationConsumer application.InvocationConsumer
}

type RepositoryManager struct {
	Machine application.MachineRepository
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

	log.Info().Msgf("Ping MongoDB successful for database %s and collection %s", config.Database, config.Collection)

	collection := client.Database(config.Database).Collection(config.Collection)

	return repository.NewMachineRepository(repository.NewMachineDataAccessMapper(), repository.NewMachineMongoRepository(collection)), nil
}

func setupInvocationConsumer(ctx context.Context, repositoryManager *RepositoryManager) (application.InvocationConsumer, error) {
	// Application Service
	applicationService := application.NewInvocationEventHandler(
		application.NewInvocationEventHandlerRepositoryManager(repositoryManager.Machine),
	)

	// Config
	kafkaConfig, err := config.SetupInvocationConsumerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	// Setup Kafka consumer
	deserializer := messaging.NewInvocationMessageDeserializer()
	processor := messaging.NewInvocationMessageProcessor(applicationService)

	confluentConsumer, err := kafka.NewConfluentKafkaConsumer(ctx, kafkaConfig, deserializer, processor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	consumer := messaging.NewInvocationConsumer(confluentConsumer)

	log.Info().Msg("Invocation consumer successfully initialized")

	return consumer, nil
}

// Consumer

func startConsumer(consumer application.InvocationConsumer, shutdown <-chan os.Signal, timeout time.Duration) {
	// Start consumer in a goroutine
	consumerErr := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting invocation consumer...")

		consumer.PollAndProcessMessages()
	}()

	// Wait for shutdown signal
	select {
	case err := <-consumerErr:
		log.Fatal().Msgf("consumer failed: %v", err)
	case sig := <-shutdown:
		log.Info().Msgf("Received signal: %s, shutting down...", sig)

		gracefulShutdown(timeout)
	}
}

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}

func gracefulShutdown(timeout time.Duration) {
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
