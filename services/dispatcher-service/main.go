package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize context
	ctx := context.Background()

	// Load environment variables
	log.Info().Msg("Loading environment variables...")
	loadEnv()

	// Setup dependencies
	log.Info().Msg("Setting up dependencies...")
	dependencies, err := setupDependencies(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup dependencies")
	}

	// Signal handling for graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Start the invocation consumer in a separate goroutine
	log.Info().Msg("Starting invocation consumer...")
	go func() {
		dependencies.invocationConsumer.PollAndProcessMessages()
	}()

	log.Info().Msg("Dispatcher service is running...")
	<-shutdownChan

	log.Info().Msg("Shutdown signal received, initiating graceful shutdown...")
}
