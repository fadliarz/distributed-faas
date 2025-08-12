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

func loadEnv(config *Config) {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found")
	}

	if timeout := os.Getenv("SHUTDOWN_TIMEOUT"); timeout != "" {
		if parsed, err := time.ParseDuration(timeout); err == nil {
			config.ShutdownTimeout = parsed
		}
	}

	if config.ShutdownTimeout == 0 {
		config.ShutdownTimeout = 30 * time.Second
	}
}

// Dependencies

type Dependencies struct {
	consumer application.BillingCalculationConsumer
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	// Load configurations
	mongoConfig, err := config.NewBillingCalculatorMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load mongo config: %w", err)
	}

	kafkaConfig, err := config.NewBillingCalculatorKafkaConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kafka config: %w", err)
	}

	// Setup MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoConfig.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	database := client.Database(mongoConfig.Database)
	chargeCollection := database.Collection(mongoConfig.ChargeCollection)
	billingCollection := database.Collection(mongoConfig.BillingCollection)

	// Setup repositories
	chargeMapper := repository.NewChargeDataAccessMapper()
	chargeMongoRepo := repository.NewChargeMongoRepository(chargeCollection)
	chargeRepository := repository.NewChargeRepository(chargeMapper, chargeMongoRepo)

	billingMapper := repository.NewBillingDataAccessMapper()
	billingMongoRepo := repository.NewBillingMongoRepository(billingCollection)
	billingRepository := repository.NewBillingRepository(billingMapper, billingMongoRepo)

	repositoryManager := application.NewBillingCalculatorApplicationServiceRepositoryManager(
		chargeRepository,
		billingRepository,
	)

	// Setup domain services
	domainService := domain.NewBillingCalculatorDomainService()

	// Setup application services
	dataMapper := application.NewBillingCalculatorDataMapper()
	applicationService := application.NewBillingCalculatorApplicationService(
		dataMapper,
		domainService,
		repositoryManager,
	)

	// Setup event handler
	eventHandler := application.NewBillingCalculationEventHandler(applicationService)

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

	return &Dependencies{
		consumer: consumer,
	}, nil
}

// Helpers

func setupShutdownHandler() <-chan struct{} {
	shutdown := make(chan struct{})
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigchan
		log.Info().Msg("Received shutdown signal")
		close(shutdown)
	}()

	return shutdown
}
