package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/services/charge-service/config"
	"github.com/fadliarz/distributed-faas/services/charge-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/charge-service/domain/domain-core"
	"github.com/fadliarz/distributed-faas/services/charge-service/infrastructure/messaging"
	"github.com/fadliarz/distributed-faas/services/charge-service/rpc"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Config

type Config struct {
	Port            string
	ShutdownTimeout time.Duration
}

// Env

func loadEnv(config *Config) {
	if err := godotenv.Load(); err != nil {
		log.Warn().Msg("No .env file found, using default configuration")
	}

	if port := os.Getenv("PORT"); port != "" {
		config.Port = ":" + port

		log.Info().Msgf("Using port from .env: %s", config.Port)
	} else {
		log.Info().Msgf("Using default port: %s", config.Port)
	}
}

// Dependencies

type Dependencies struct {
	handler *application.CommandHandler
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	chargeAggregator, err := setupChargeAggregator(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup charge aggregator: %w", err)
	}

	return &Dependencies{
		handler: setupCommandHandler(chargeAggregator),
	}, nil
}

func setupChargeAggregator(ctx context.Context) (application.ChargeAggregator, error) {
	producer, err := setupChargeProducer(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup charge producer: %w", err)
	}

	chargeAggregator := domain.NewChargeAggregator(
		domain.NewChargeDomainService(),
		producer,
		5*time.Second, // aggregation duration
	)

	// Start the aggregator
	go chargeAggregator.Start(ctx)

	return chargeAggregator, nil
}

func setupChargeProducer(ctx context.Context) (application.ChargeProducer, error) {
	producerConfig, err := config.NewChargeProducerConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load charge producer config: %w", err)
	}

	producer, err := messaging.NewChargeProducer(ctx, producerConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create charge producer: %w", err)
	}

	return producer, nil
}

func setupCommandHandler(aggregator application.ChargeAggregator) *application.CommandHandler {
	applicationService := application.NewChargeApplicationService(
		application.NewChargeDataMapper(),
		aggregator,
	)

	return application.NewCommandHandler(
		applicationService,
	)
}

// Servers

func setupGRPCServer(config *Config, dependencies *Dependencies) (*grpc.Server, net.Listener, error) {
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		return nil, nil, err
	}

	log.Info().Msgf("gRPC server listening on %s", config.Port)

	grpcServer := grpc.NewServer()
	chargeServer := rpc.NewChargeServer(dependencies.handler)
	chargeServer.Register(grpcServer)
	reflection.Register(grpcServer)

	return grpcServer, lis, nil
}

func startServer(server *grpc.Server, listener net.Listener, shutdown <-chan os.Signal, timeout time.Duration) {
	// Start server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Info().Msg("Starting gRPC server...")

		if err := server.Serve(listener); err != nil {
			serverErr <- err
		}
	}()

	// Wait for shutdown signal or server error
	select {
	case err := <-serverErr:
		log.Fatal().Msgf("server failed: %v", err)
	case sig := <-shutdown:
		log.Info().Msgf("Received signal: %s, shutting down...", sig)

		gracefulShutdown(server, timeout)
	}
}

func setupShutdownHandler() <-chan os.Signal {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	return sigChan
}

func gracefulShutdown(server *grpc.Server, timeout time.Duration) {
	log.Info().Msg("Gracefully stopping server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)

	defer cancel()

	done := make(chan struct{})
	go func() {
		server.GracefulStop()

		close(done)
	}()

	select {
	case <-shutdownCtx.Done():
		log.Warn().Msg("Shutdown timeout exceeded, forcing stop")

		server.Stop()
	case <-done:
		log.Info().Msg("Server stopped gracefully")
	}
}
