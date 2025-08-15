package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/services/accumulator-service/config"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/infrastructure/messaging"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/infrastructure/repository"
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

	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		if parsed, err := time.ParseDuration(timeout); err == nil {
			config.ShutdownTimeout = parsed
			log.Info().Msgf("Using shutdown timeout from .env: %s", config.ShutdownTimeout)
		}
	}

	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 30 * time.Second
		log.Info().Msgf("Using default shutdown timeout: %s", config.ShutdownTimeout)
	}
}

// Dependencies

type Dependencies struct {
	consumer application.ChargeConsumer
}

type RepositoryManager struct {
	Charge application.ChargeRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	consumer, err := setupKafkaConsumer(repositoryManager)
	if err != nil {
		return nil, fmt.Errorf("failed to setup kafka consumer: %w", err)
	}

	return &Dependencies{
		consumer: consumer,
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	chargeRepository, err := setupChargeRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup charge repository: %w", err)
	}

	return &RepositoryManager{
		Charge: chargeRepository,
	}, nil
}

func setupChargeRepository(ctx context.Context) (application.ChargeRepository, error) {
	// Config
	config, err := config.NewChargeMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load accumulator mongo config: %w", err)
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

	// Repository layers
	mongoRepository := repository.NewChargeMongoRepository(collection)
	dataAccessMapper := repository.NewChargeDataAccessMapper()
	repo := repository.NewChargeRepository(dataAccessMapper, mongoRepository)

	return repo, nil
}

func setupKafkaConsumer(repositoryManager *RepositoryManager) (application.ChargeConsumer, error) {
	// Application Service
	applicationService := application.NewChargeApplicationService(repositoryManager.Charge)
	eventHandler := application.NewChargeEventHandler(applicationService)

	// Config
	kafkaConfig, err := config.NewChargeKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	// Setup Kafka consumer
	consumer, err := messaging.NewChargeConsumer(kafkaConfig, eventHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	log.Info().Msg("Charge consumer successfully initialized")

	return consumer, nil
}

// Consumer

func startConsumer(consumer application.ChargeConsumer, shutdown <-chan os.Signal, timeout time.Duration) {
	// Start consumer in a goroutine
	consumerErr := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting charge consumer...")

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
