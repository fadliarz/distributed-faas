package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/services/billing-cron-service/config"
	"github.com/fadliarz/distributed-faas/services/billing-cron-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-cron-service/infrastructure/repository"
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
	Handler       *application.BillingCronHandler
}

type ConfigManager struct {
	BillingCron *config.BillingCronConfig
}

type RepositoryManager struct {
	Cron application.CronRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	billingCronConfig, err := config.NewBillingCronConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load billing cron config: %w", err)
	}

	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	return &Dependencies{
		ConfigManager: &ConfigManager{
			BillingCron: billingCronConfig,
		},
		Handler: setupBillingCronHandler(repositoryManager),
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	cronRepository, err := setupCronRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup cron store repository: %w", err)
	}

	return &RepositoryManager{
		Cron: cronRepository,
	}, nil
}

func setupCronRepository(ctx context.Context) (application.CronRepository, error) {
	// Config
	config, err := config.NewCronMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load cron store mongo config: %w", err)
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

	return repository.NewCronRepository(repository.NewCronMongoRepository(collection)), nil
}

func setupBillingCronHandler(repositoryManager *RepositoryManager) *application.BillingCronHandler {
	// Application Service
	applicationService := application.NewBillingCronApplicationService(
		*application.NewBillingCronApplicationServiceRepositoryManager(repositoryManager.Cron),
	)

	return application.NewBillingCronHandler(applicationService)
}

// Service

func startBillingCronService(handler *application.BillingCronHandler, config *config.BillingCronConfig, shutdown <-chan os.Signal, timeout time.Duration, ctx context.Context) {
	// Start service in a goroutine
	serviceErr := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting billing cron service...")

		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("Context cancelled, shutting down")

				return
			default:
				err := handler.UpdateLastBilled(ctx)
				if err != nil {
					log.Error().Err(err).Msg("An error occurred during LastBilled update")
				} else {
					log.Info().Msg("LastBilled update completed successfully")
				}

				time.Sleep(time.Duration(config.CronIntervalInSec) * time.Second)
			}
		}
	}()

	// Wait for shutdown signal
	select {
	case err := <-serviceErr:
		log.Fatal().Msgf("service failed: %v", err)
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
	log.Info().Msg("Gracefully stopping service...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	// Wait for the timeout
	<-shutdownCtx.Done()

	log.Info().Msg("Service stopped gracefully")
}
