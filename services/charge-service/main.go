package main

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	config := &Config{
		Port:            ":50050",
		ShutdownTimeout: 30 * time.Second,
	}

	loadEnv(config)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dependencies, err := setupDependencies(ctx)
	if err != nil {
		log.Fatal().Msgf("failed to setup dependencies: %v", err)
	}

	server, listener, err := setupGRPCServer(config, dependencies)
	if err != nil {
		log.Fatal().Msgf("failed to setup gRPC server: %v", err)
	}

	shutdown := setupShutdownHandler()

	startServer(server, listener, shutdown, config.ShutdownTimeout)
}
