package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/config"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/infrastructure/messaging"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/infrastructure/repository"
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
	consumer application.BillingCalculationConsumer
}

type RepositoryManager struct {
	Charge  application.ChargeRepository
	Billing application.BillingRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	consumer, err := setupConsumer(ctx, repositoryManager)
	if err != nil {
		return nil, fmt.Errorf("failed to setup consumer: %w", err)
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

	billingRepository, err := setupBillingRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup billing repository: %w", err)
	}

	return &RepositoryManager{
		Charge:  chargeRepository,
		Billing: billingRepository,
	}, nil
}

func setupChargeRepository(ctx context.Context) (application.ChargeRepository, error) {
	// Config
	config, err := config.NewChargeMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load charge mongo config: %w", err)
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

	return repository.NewChargeRepository(repository.NewChargeDataAccessMapper(), repository.NewChargeMongoRepository(collection)), nil
}

func setupBillingRepository(ctx context.Context) (application.BillingRepository, error) {
	// Config
	config, err := config.NewBillingMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load billing mongo config: %w", err)
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

	return repository.NewBillingRepository(repository.NewBillingDataAccessMapper(), repository.NewBillingMongoRepository(collection)), nil
}

func setupConsumer(ctx context.Context, repositoryManager *RepositoryManager) (application.BillingCalculationConsumer, error) {
	// Application Service
	applicationService := application.NewBillingCalculatorApplicationService(
		application.NewBillingCalculatorDataMapper(),
		domain.NewBillingCalculatorDomainService(),
		application.NewBillingCalculatorApplicationServiceRepositoryManager(
			repositoryManager.Charge,
			repositoryManager.Billing,
		),
	)

	// Event Handler
	eventHandler := application.NewBillingCalculationEventHandler(applicationService)

	// Kafka Config
	kafkaConfig, err := config.NewCronKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	// Setup Kafka consumer
	kafkaConsumerConfig := &kafka.ConsumerConfig{
		Basic: &kafka.ConsumerBasicConfig{
			BootstrapServers: kafkaConfig.BootstrapServers,
			Topic:            kafkaConfig.Topic,
			GroupID:          kafkaConfig.GroupID,
			PollTimeout:      kafkaConfig.PollTimeout,
		},
		Processing: &kafka.ConsumerProcessingConfig{
			NumWorkers: kafkaConfig.NumWorkers,
		},
	}

	deserializer := messaging.NewBillingCalculationMessageDeserializer()
	processor := messaging.NewBillingCalculationMessageProcessor(eventHandler)

	kafkaConsumer, err := kafka.NewConfluentKafkaConsumer(ctx, kafkaConsumerConfig, deserializer, processor)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	consumer := messaging.NewBillingCalculationConsumer(kafkaConsumer)

	log.Info().Msg("All dependencies successfully initialized")

	return consumer, nil
}

// Consumer

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}
