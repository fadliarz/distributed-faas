package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	config := &Config{
		ShutdownTimeout: 30 * time.Second,
	}

	loadEnv(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dependencies, err := setupDependencies(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to setup dependencies: %v", err)
	}

	shutdown := setupShutdownHandler()

	go func() {
		dependencies.consumer.PollAndProcessMessages()
	}()

	<-shutdown
}
