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

	consumer, err := setupKafkaConsumer(ctx, repositoryManager)
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
	config, err := config.NewAccumulatorMongoConfig()
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

func setupKafkaConsumer(ctx context.Context, repositoryManager *RepositoryManager) (application.ChargeConsumer, error) {
	// Config
	kafkaConfig, err := config.NewAccumulatorKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	// Application service and handler
	chargeService := application.NewChargeApplicationService(repositoryManager.Charge)
	eventHandler := application.NewChargeEventHandler(chargeService)

	// Custom Kafka consumer
	consumer, err := messaging.NewChargeConsumer(kafkaConfig, eventHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka consumer: %w", err)
	}

	return consumer, nil
}

// Shutdown

func setupShutdownHandler() <-chan os.Signal {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	return shutdown
}
