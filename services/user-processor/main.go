package main

import (
	"context"

	"github.com/rs/zerolog/log"
)

func main() {
	loadEnv()

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
