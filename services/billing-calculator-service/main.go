package main

import (
	"context"

	"github.com/rs/zerolog/log"
)

func main() {
	config := &Config{}

	loadEnv(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dependencies, err := setupDependencies(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup dependencies")
	}

	log.Info().Msg("Starting billing calculator service")

	shutdown := setupShutdownHandler()

	// Start the consumer
	go func() {
		log.Info().Msg("Starting billing calculation consumer")
		dependencies.consumer.PollAndProcessMessages()
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Info().Msg("Shutting down billing calculator service")
}
