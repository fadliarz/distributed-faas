package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/compose"
	"golang.org/x/sync/errgroup"
)

type ContainerManager struct {
	ctx    context.Context
	config *TestConfig

	ConnectionStrings *ConnectionStrings

	Composes   *Composes
	Containers *Containers
}

type Composes struct {
	Mongo     *compose.DockerCompose
	Zookeeper *compose.DockerCompose
	Kafka     *compose.DockerCompose
	Services  *compose.DockerCompose
}

type Containers struct {
	Mongo       testcontainers.Container
	Kafka       testcontainers.Container
	UserService testcontainers.Container
}

type ConnectionStrings struct {
	Mongo string
}

func NewContainerManager(ctx context.Context, config *TestConfig) *ContainerManager {
	return &ContainerManager{
		ctx:               ctx,
		config:            config,
		ConnectionStrings: &ConnectionStrings{},
		Composes:          &Composes{},
		Containers:        &Containers{},
	}
}

func (cm *ContainerManager) SetupContainers() error {
	if err := cm.setupComposes(); err != nil {
		return fmt.Errorf("failed to setup infrastructure composes: %w", err)
	}

	if err := cm.startContainers(); err != nil {
		return fmt.Errorf("failed to start infrastructure containers: %w", err)
	}

	if err := cm.setupConnectionStrings(); err != nil {
		return fmt.Errorf("failed to setup connection strings: %w", err)
	}

	return nil
}

func (cm *ContainerManager) setupComposes() error {
	var err error

	// Mongo
	cm.Composes.Mongo, err = compose.NewDockerComposeWith(
		compose.StackIdentifier(cm.config.ComposeConfig.ProjectID),
		compose.WithStackFiles(cm.config.ComposePaths.Common, cm.config.ComposePaths.Mongo),
	)
	if err != nil {
		return fmt.Errorf("failed to create MongoDB Docker Compose stack: %w", err)
	}

	// Zookeeper
	cm.Composes.Zookeeper, err = compose.NewDockerComposeWith(
		compose.StackIdentifier(cm.config.ComposeConfig.ProjectID),
		compose.WithStackFiles(cm.config.ComposePaths.Common, cm.config.ComposePaths.Zookeeper),
	)
	if err != nil {
		return fmt.Errorf("failed to create Zookeeper Docker Compose stack: %w", err)
	}

	// Kafka
	cm.Composes.Kafka, err = compose.NewDockerComposeWith(
		compose.StackIdentifier(cm.config.ComposeConfig.ProjectID),
		compose.WithStackFiles(cm.config.ComposePaths.Common, cm.config.ComposePaths.Kafka),
		compose.WithProfiles(cm.config.ComposeConfig.Profile),
	)
	if err != nil {
		return fmt.Errorf("failed to create Kafka Docker Compose stack: %w", err)
	}

	// Services
	cm.Composes.Services, err = compose.NewDockerComposeWith(
		compose.StackIdentifier(cm.config.ComposeConfig.ProjectID),
		compose.WithStackFiles(cm.config.ComposePaths.Common, cm.config.ComposePaths.Services),
		compose.WithProfiles(cm.config.ComposeConfig.Profile),
	)
	if err != nil {
		return fmt.Errorf("failed to create Services Docker Compose stack: %w", err)
	}

	return nil
}

func (cm *ContainerManager) startContainers() error {
	g := new(errgroup.Group)

	g.Go(func() error {
		// Start Mongo
		log.Info().Msg("Starting MongoDB compose...")
		err := cm.Composes.Mongo.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start MongoDB Docker Compose stack: %w", err)
		}

		cm.Containers.Mongo, err = cm.Composes.Mongo.ServiceContainer(cm.ctx, cm.config.ContainerNames.Mongo)
		if err != nil {
			return fmt.Errorf("failed to get MongoDB container: %w", err)
		}

		// Start Services
		log.Info().Msg("Starting Services compose...")
		err = cm.Composes.Services.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start Services Docker Compose stack: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		// Start Zookeeper
		log.Info().Msg("Starting Zookeeper compose...")
		err := cm.Composes.Zookeeper.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start Zookeeper Docker Compose stack: %w", err)
		}

		// Start Kafka
		log.Info().Msg("Starting Kafka compose...")
		err = cm.Composes.Kafka.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start Kafka Docker Compose stack: %w", err)
		}

		cm.Containers.Kafka, err = cm.Composes.Kafka.ServiceContainer(cm.ctx, cm.config.ContainerNames.Kafka)
		if err != nil {
			return fmt.Errorf("failed to get Kafka container: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to wait for container startup: %w", err)
	}

	return nil
}

func (cm *ContainerManager) setupConnectionStrings() error {
	host, err := cm.Containers.Mongo.Host(cm.ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB container host: %w", err)
	}

	port, err := cm.Containers.Mongo.MappedPort(cm.ctx, "27017")
	if err != nil {
		return fmt.Errorf("failed to get MongoDB container port: %w", err)
	}

	cm.ConnectionStrings.Mongo = "mongodb://" + cm.config.MongoConfig.Username + ":" + cm.config.MongoConfig.Password + "@" + host + ":" + port.Port() + fmt.Sprintf("/?replicaSet=%s&directConnection=true", cm.config.MongoConfig.ReplicaSet)

	return nil
}

func (cm *ContainerManager) Down() error {
	g := new(errgroup.Group)

	// Stop Services
	if cm.Composes.Services != nil {
		g.Go(func() error {
			return cm.Composes.Services.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
		})
	}

	// Stop Kafka
	if cm.Composes.Kafka != nil {
		g.Go(func() error {
			return cm.Composes.Kafka.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
		})
	}

	// Stop Zookeeper
	if cm.Composes.Zookeeper != nil {
		g.Go(func() error {
			return cm.Composes.Zookeeper.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
		})
	}

	// Stop Mongo
	if cm.Composes.Mongo != nil {
		g.Go(func() error {
			return cm.Composes.Mongo.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveImagesLocal)
		})
	}

	return g.Wait()
}
