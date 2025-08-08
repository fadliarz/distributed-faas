package main

import (
	"context"
	"time"

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
		for {
			select {
			case <-ctx.Done():
				log.Info().Msg("context cancelled, shutting down")

				return
			default:
				err := dependencies.Handler.RetryInvocations(ctx, dependencies.ConfigManager.Retry.ThresholdInSec)

				if err != nil {
					log.Error().Err(err).Msg("an error occured")
				} else {
					log.Info().Msg("retry invocations completed successfully")
				}

				time.Sleep(time.Duration(dependencies.ConfigManager.Retry.RetryIntervalInSec) * time.Second)
			}
		}
	}()

	<-shutdown
}
