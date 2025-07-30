package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fadliarz/distributed-faas/services/machine/config"
	"github.com/fadliarz/distributed-faas/services/machine/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"github.com/fadliarz/distributed-faas/services/machine/infrastructure/repository"
	"github.com/fadliarz/distributed-faas/services/machine/infrastructure/storage"
	"github.com/fadliarz/distributed-faas/services/machine/rpc"
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
	Checkpoint application.CheckpointRepository
}

func setupDependencies(ctx context.Context) (*Dependencies, error) {
	// Repository Manager
	repositoryManager, err := setupRepositoryManager(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to setup repository manager: %w", err)
	}

	// Command Handler
	commandHandler, err := setupCommandHandler(ctx, repositoryManager)
	if err != nil {
		return nil, fmt.Errorf("failed to setup command handler: %w", err)
	}

	return &Dependencies{
		handler: commandHandler,
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
	cfg, err := config.NewCheckpointMongoConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load checkpoint MongoDB config: %v", err)
	}

	// Collection
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	collection := client.Database(cfg.Database).Collection(cfg.Collection)

	return repository.NewCheckpointRepositoryImpl(repository.NewCheckpointDataAccessMapperImpl(), repository.NewCheckpointMongoRepository(collection)), nil
}

func setupCommandHandler(ctx context.Context, repositoryManager *RepositoryManager) (*application.CommandHandler, error) {
	// Application Service
	applicationService := application.NewMachineApplicationService(
		application.NewMachineDataMapper(),
		domain.NewMachineDomainService(),
		application.NewMachineApplicationServiceRepositoryManager(repositoryManager.Checkpoint),
	)

	// Config
	cloudflareConfig, err := config.NewOutputCloudflareConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load Cloudflare config: %v", err)
	}

	// Client
	s3Client, err := storage.NewOutputS3Client(ctx, cloudflareConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 client: %v", err)
	}

	return application.NewCommandHandler(
		applicationService,
		application.NewCommandHandlerConfig(*cloudflareConfig),
		application.NewCommandHandlerClient(s3Client),
	), nil
}

// Servers

func setupGRPCServer(ctx context.Context, config *Config, dependencies *Dependencies) (*grpc.Server, net.Listener, error) {
	// Create a TCP listener on the specified port
	lis, err := net.Listen("tcp", config.Port)
	if err != nil {
		return nil, nil, err
	}

	log.Info().Msgf("gRPC server listening on %s", config.Port)

	// Create gRPC server and register the function server
	grpcServer := grpc.NewServer()
	functionServer := rpc.NewMachineServer(ctx, dependencies.handler)
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
