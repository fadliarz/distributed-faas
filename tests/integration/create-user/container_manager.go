package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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
	Debezium    testcontainers.Container
	UserService testcontainers.Container
}

type ConnectionStrings struct {
	Mongo       string
	Debezium    string
	UserService string
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

	if err := cm.setupDebeziumConnectors(); err != nil {
		return fmt.Errorf("failed to setup Debezium connector: %w", err)
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

		cm.Containers.UserService, err = cm.Composes.Services.ServiceContainer(cm.ctx, cm.config.ContainerNames.UserService)
		if err != nil {
			return fmt.Errorf("failed to get User Service container: %w", err)
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

		// Get Debezium container
		cm.Containers.Debezium, err = cm.Composes.Kafka.ServiceContainer(cm.ctx, cm.config.ContainerNames.Debezium)
		if err != nil {
			return fmt.Errorf("failed to get Debezium container: %w", err)
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to wait for container startup: %w", err)
	}

	return nil
}

func (cm *ContainerManager) setupConnectionStrings() error {
	// MongoDB connection string
	host, err := cm.Containers.Mongo.Host(cm.ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB container host: %w", err)
	}

	port, err := cm.Containers.Mongo.MappedPort(cm.ctx, "27017")
	if err != nil {
		return fmt.Errorf("failed to get MongoDB container port: %w", err)
	}

	cm.ConnectionStrings.Mongo = "mongodb://" + cm.config.MongoConfig.Username + ":" + cm.config.MongoConfig.Password + "@" + host + ":" + port.Port() + fmt.Sprintf("/?replicaSet=%s&directConnection=true", cm.config.MongoConfig.ReplicaSet)

	// Debezium connection string
	debeziumHost, err := cm.Containers.Debezium.Host(cm.ctx)
	if err != nil {
		return fmt.Errorf("failed to get Debezium host: %w", err)
	}

	debeziumPort, err := cm.Containers.Debezium.MappedPort(cm.ctx, "8083")
	if err != nil {
		return fmt.Errorf("failed to get Debezium port: %w", err)
	}

	cm.ConnectionStrings.Debezium = fmt.Sprintf("http://%s:%s", debeziumHost, debeziumPort.Port())

	// User Service
	// Invocation Service
	host, err = cm.Containers.UserService.Host(cm.ctx)
	if err != nil {
		return fmt.Errorf("failed to get User Service container host: %w", err)
	}

	port, err = cm.Containers.UserService.MappedPort(cm.ctx, "50050")
	if err != nil {
		return fmt.Errorf("failed to get User Service container port: %w", err)
	}

	cm.ConnectionStrings.UserService = host + ":" + port.Port()

	return nil
}

func (cm *ContainerManager) setupDebeziumConnectors() error {
	err := cm.setupUserDebeziumConnector()
	if err != nil {
		return fmt.Errorf("failed to setup user Debezium connector: %w", err)
	}

	return nil
}

func (cm *ContainerManager) setupUserDebeziumConnector() error {
	// Set up the Debezium connector configuration
	configJSON := fmt.Sprintf(`{
		"name": "%s",
		"config": {
			"connector.class": "io.debezium.connector.mongodb.MongoDbConnector",
			"mongodb.connection.string": "%s",
			"topic.prefix": "cdc",
			"database.include.list": "%s",
			"collection.include.list": "%s",

			"key.converter": "org.apache.kafka.connect.json.JsonConverter",
			"key.converter.schemas.enable": false,
			"value.converter": "org.apache.kafka.connect.json.JsonConverter",
			"value.converter.schemas.enable": false,

			"transforms": "filter,unwrap",

			"transforms.filter.type": "io.debezium.transforms.Filter",
			"transforms.filter.language": "jsr223.groovy",
			"transforms.filter.condition": "value.op == 'c'",

			"transforms.unwrap.type": "io.debezium.connector.mongodb.transforms.ExtractNewDocumentState"
		}
	}`,
		cm.config.DebeziumConfig.UserConnectorName,
		cm.config.GetMongoConnectionString(),
		cm.config.MongoConfig.UserDatabase,
		fmt.Sprintf("%s.%s", cm.config.MongoConfig.UserDatabase, cm.config.MongoConfig.UserCollection),
	)

	fmt.Print(configJSON)

	endpoint := fmt.Sprintf("%s/connectors", cm.ConnectionStrings.Debezium)

	for attempt := 1; attempt <= cm.config.DebeziumConfig.MaxRetries; attempt++ {
		log.Debug().Msgf("[%s] Attempting to create Debezium connector (attempt %d/%d)", cm.config.DebeziumConfig.UserConnectorName, attempt, cm.config.DebeziumConfig.MaxRetries)

		// Create HTTP request
		req, err := http.NewRequest("POST", endpoint, strings.NewReader(configJSON))
		if err != nil {
			return fmt.Errorf("[%s] failed to create HTTP request: %w", cm.config.DebeziumConfig.UserConnectorName, err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create HTTP client with timeout
		client := &http.Client{Timeout: 10 * time.Second}

		// Send HTTP request
		res, err := client.Do(req)

		// Handle response
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to create Debezium connector", cm.config.DebeziumConfig.UserConnectorName)
		}

		if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Debug().Msgf("[%s] Debezium connector created successfully on attempt %d", cm.config.DebeziumConfig.UserConnectorName, attempt)

			res.Body.Close()
			return cm.waitForUserConnectorReady(fmt.Sprintf("%s/%s/status", endpoint, cm.config.DebeziumConfig.UserConnectorName))
		}

		if res != nil {
			body := ""
			if b, readErr := io.ReadAll(res.Body); readErr == nil {
				body = string(b)
			}
			res.Body.Close()

			log.Error().Msgf("[%s] Failed to create Debezium connector, status code: %d, response body: %s", cm.config.DebeziumConfig.UserConnectorName, res.StatusCode, body)
		}

		if attempt < cm.config.DebeziumConfig.MaxRetries {
			time.Sleep(cm.config.DebeziumConfig.RetryInterval)
		}
	}

	return fmt.Errorf("[%s] failed to create Debezium connector after %d attempts", cm.config.DebeziumConfig.UserConnectorName, cm.config.DebeziumConfig.MaxRetries)
}

func (cm *ContainerManager) waitForUserConnectorReady(endpoint string) error {
	deadline := time.Now().Add(cm.config.DebeziumConfig.ReadyTimeout)

	for time.Now().Before(deadline) {
		log.Debug().Msgf("[%s] Checking Debezium connector status", cm.config.DebeziumConfig.UserConnectorName)

		res, err := http.Get(endpoint)

		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to get Debezium connector status", cm.config.DebeziumConfig.UserConnectorName)
			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
			continue
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to read Debezium connector response body", cm.config.DebeziumConfig.UserConnectorName)
			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
			continue
		}

		if strings.Contains(string(body), `"state":"RUNNING"`) {
			log.Debug().Msgf("[%s] Debezium connector is ready", cm.config.DebeziumConfig.UserConnectorName)
			return nil
		}

		time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
	}

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
