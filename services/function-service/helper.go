package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/services/function-service/config"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
	"github.com/fadliarz/distributed-faas/services/function-service/infrastructure/repository"
	"github.com/fadliarz/distributed-faas/services/function-service/rpc"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type RepositoryManager struct {
	Function application.FunctionRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	return &Dependencies{
		handler: setupCommandHandler(repositoryManager),
	}, nil
}

func setupRepositoryManager(ctx context.Context) (*RepositoryManager, error) {
	functionRepository, err := setupFunctionRepository(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup function repository: %w", err)
	}

	return &RepositoryManager{
		Function: functionRepository,
	}, nil
}

func setupFunctionRepository(ctx context.Context) (application.FunctionRepository, error) {
	// Config
	config, err := config.NewFunctionMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load function mongo config: %w", err)
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

	return repository.NewFunctionRepository(repository.NewFunctionDataAccessMapper(), repository.NewFunctionMongoRepository(collection)), nil
}

func setupCommandHandler(repositoryManager *RepositoryManager) *application.CommandHandler {
	// Application Service
	applicationService := application.NewFunctionApplicationService(
		application.NewFunctionDataMapper(),
		domain.NewFunctionDomainService(),
		application.NewFunctionApplicationServiceRepositoryManager(repositoryManager.Function),
	)

	return application.NewCommandHandler(
		applicationService,
	)
}

// Servers

func setupGRPCServer(config *Config, dependencies *Dependencies) (*grpc.Server, net.Listener, error) {
	// Create a TCP listener on the specified port
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		return nil, nil, err
	}

	log.Info().Msgf("gRPC server listening on %s", config.Port)

	// Create gRPC server and register the function server
	grpcServer := grpc.NewServer()
	functionServer := rpc.NewFunctionServer(dependencies.handler)
	functionServer.Register(grpcServer)
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
