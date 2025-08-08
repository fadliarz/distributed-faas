package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/services/retry-service/config"
	"github.com/fadliarz/distributed-faas/services/retry-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/retry-service/infrastructure/repository"
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
	ConfigManager *ConfigManager
	Handler       *application.RetryHandler
}

type ConfigManager struct {
	Retry *config.RetryConfig
}

type RepositoryManager struct {
	Checkpoint application.CheckpointRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	retryConfig, err := config.NewRetryConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load retry config: %w", err)
	}

	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	return &Dependencies{
		ConfigManager: &ConfigManager{
			Retry: retryConfig,
		},
		Handler: setupRetryHandler(repositoryManager),
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	checkpointRepository, err := setupCheckpointRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup checkpoint repository: %w", err)
	}

	return &RepositoryManager{
		Checkpoint: checkpointRepository,
	}, nil
}

func setupCheckpointRepository(ctx context.Context) (application.CheckpointRepository, error) {
	// Config
	config, err := config.NewCheckpointMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load checkpoint mongo config: %w", err)
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

	return repository.NewCheckpointRepository(repository.NewCheckpointMongoRepository(collection)), nil
}

func setupRetryHandler(repositoryManager *RepositoryManager) *application.RetryHandler {
	// Confluent Consumer
	applicationService := application.NewRetryApplicationService(*application.NewRetryApplicationServiceRepositoryManager(repositoryManager.Checkpoint))

	return application.NewRetryHandler(applicationService)
}

// Servers

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}
