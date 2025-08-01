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
	Mongo             testcontainers.Container
	Kafka             testcontainers.Container
	Debezium          testcontainers.Container
	FunctionService   testcontainers.Container
	InvocationService testcontainers.Container
}

type ConnectionStrings struct {
	Mongo    string
	Debezium string
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
	)
	if err != nil {
		return fmt.Errorf("failed to create Kafka Docker Compose stack: %w", err)
	}

	// Services
	cm.Composes.Services, err = compose.NewDockerComposeWith(
		compose.StackIdentifier(cm.config.ComposeConfig.ProjectID),
		compose.WithStackFiles(cm.config.ComposePaths.Common, cm.config.ComposePaths.Services),
	)
	if err != nil {
		return fmt.Errorf("failed to create Services Docker Compose stack: %w", err)
	}

	return nil
}

func (cm *ContainerManager) startContainers() error {
	g := new(errgroup.Group)

	g.Go(func() error {
		err := cm.Composes.Mongo.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start MongoDB Docker Compose stack: %w", err)
		}

		cm.Containers.Mongo, err = cm.Composes.Mongo.ServiceContainer(cm.ctx, cm.config.ContainerNames.Mongo)
		if err != nil {
			return fmt.Errorf("failed to get MongoDB container: %w", err)
		}

		// Services
		err = cm.Composes.Services.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start Services Docker Compose stack: %w", err)
		}

		cm.Containers.FunctionService, err = cm.Composes.Services.ServiceContainer(cm.ctx, cm.config.ContainerNames.FunctionService)
		if err != nil {
			return fmt.Errorf("failed to get Function Service container: %w", err)
		}

		return nil
	})

	g.Go(func() error {
		// Zookeeper
		err := cm.Composes.Zookeeper.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start Zookeeper Docker Compose stack: %w", err)
		}

		// Kafka
		err = cm.Composes.Kafka.Up(cm.ctx, compose.Wait(true))
		if err != nil {
			return fmt.Errorf("failed to start Kafka Docker Compose stack: %w", err)
		}

		cm.Containers.Kafka, err = cm.Composes.Kafka.ServiceContainer(cm.ctx, cm.config.ContainerNames.Kafka)
		if err != nil {
			return fmt.Errorf("failed to get Kafka container: %w", err)
		}

		// Debezium
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
	// Mongo
	host, err := cm.Containers.Mongo.Host(cm.ctx)
	if err != nil {
		return fmt.Errorf("failed to get MongoDB container host: %w", err)
	}

	port, err := cm.Containers.Mongo.MappedPort(cm.ctx, "27017")
	if err != nil {
		return fmt.Errorf("failed to get MongoDB container port: %w", err)
	}

	cm.ConnectionStrings.Mongo = "mongodb://" + cm.config.MongoConfig.Username + ":" + cm.config.MongoConfig.Password + "@" + host + ":" + port.Port() + fmt.Sprintf("/?replicaSet=%s&directConnection=true", cm.config.MongoConfig.ReplicaSet)

	// Debezium
	host, err = cm.Containers.Debezium.Host(cm.ctx)
	if err != nil {
		return fmt.Errorf("failed to get Debezium host: %w", err)
	}

	port, err = cm.Containers.Debezium.MappedPort(cm.ctx, "8083")
	if err != nil {
		return fmt.Errorf("failed to get Debezium port: %w", err)
	}

	cm.ConnectionStrings.Debezium = fmt.Sprintf("http://%s:%s", host, port.Port())

	return nil
}

func (cm *ContainerManager) setupDebeziumConnectors() error {
	err := cm.setupInvocationDebeziumConnector()
	if err != nil {
		return fmt.Errorf("failed to setup invocation Debezium connector: %w", err)
	}

	err = cm.setupCheckpointDebeziumConnector()
	if err != nil {
		return fmt.Errorf("failed to setup checkpoint Debezium connector: %w", err)
	}

	err = cm.setupCheckpointToInvocationDebeziumConnector()
	if err != nil {
		return fmt.Errorf("failed to setup checkpoint to invocation Debezium connector: %w", err)
	}

	return nil
}

func (cm *ContainerManager) setupInvocationDebeziumConnector() error {
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
		cm.config.DebeziumConfig.InvocationConnectorName,
		cm.config.GetMongoConnectionString(),
		cm.config.MongoConfig.InvocationDatabase,
		fmt.Sprintf("%s.%s", cm.config.MongoConfig.InvocationDatabase, cm.config.MongoConfig.InvocationCollection),
	)

	endpoint := fmt.Sprintf("%s/connectors", cm.ConnectionStrings.Debezium)

	for attempt := 1; attempt <= cm.config.DebeziumConfig.MaxRetries; attempt++ {
		log.Debug().Msgf("[%s] Attempting to create Debezium connector (attempt %d/%d)", cm.config.DebeziumConfig.InvocationConnectorName, attempt, cm.config.DebeziumConfig.MaxRetries)

		// Create HTTP request
		req, err := http.NewRequest("POST", endpoint, strings.NewReader(configJSON))
		if err != nil {
			return fmt.Errorf("[%s] failed to create HTTP request: %w", cm.config.DebeziumConfig.InvocationConnectorName, err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create HTTP client with timeout
		client := &http.Client{Timeout: 10 * time.Second}

		// Send HTTP request
		res, err := client.Do(req)

		// Handle response
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to create Debezium connector", cm.config.DebeziumConfig.InvocationConnectorName)
		}

		if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Debug().Msgf("[%s] Debezium connector created successfully on attempt %d", cm.config.DebeziumConfig.InvocationConnectorName, attempt)

			res.Body.Close()
			return cm.waitForInvocationConnectorReady(fmt.Sprintf("%s/%s/status", endpoint, cm.config.DebeziumConfig.InvocationConnectorName))
		}

		if res != nil {
			body := ""
			if b, readErr := io.ReadAll(res.Body); readErr == nil {
				body = string(b)
			}
			res.Body.Close()

			log.Error().Msgf("[%s] Failed to create Debezium connector, status code: %d, response body: %s", cm.config.DebeziumConfig.InvocationConnectorName, res.StatusCode, body)
		}

		if attempt < cm.config.DebeziumConfig.MaxRetries {
			time.Sleep(cm.config.DebeziumConfig.RetryInterval * time.Duration(attempt))
		}
	}

	return fmt.Errorf("[%s] failed to create Debezium connector after %d attempts", cm.config.DebeziumConfig.InvocationConnectorName, cm.config.DebeziumConfig.MaxRetries)
}

func (cm *ContainerManager) waitForInvocationConnectorReady(endpoint string) error {
	deadline := time.Now().Add(cm.config.DebeziumConfig.ReadyTimeout)

	for time.Now().Before(deadline) {
		log.Debug().Msgf("[%s] Checking Debezium connector status", cm.config.DebeziumConfig.InvocationConnectorName)

		res, err := http.Get(endpoint)

		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to get Debezium connector status", cm.config.DebeziumConfig.InvocationConnectorName)

			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to read Debezium connector response body", cm.config.DebeziumConfig.InvocationConnectorName)

			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
		}

		if strings.Contains(string(body), `"state":"RUNNING"`) {
			log.Debug().Msgf("[%s] Debezium connector is ready", cm.config.DebeziumConfig.InvocationConnectorName)

			time.Sleep(5 * time.Second)
			return nil
		}

		time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
	}

	return nil
}

func (cm *ContainerManager) setupCheckpointDebeziumConnector() error {
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
			
			"transforms": "unwrap,filter",

			"transforms.unwrap.type": "io.debezium.connector.mongodb.transforms.ExtractNewDocumentState",

			"transforms.filter.type": "io.debezium.transforms.Filter",
			"transforms.filter.language": "jsr223.groovy",
			"transforms.filter.condition": "value.status == 'SUCCESS'"
		}
	}`,
		cm.config.DebeziumConfig.CheckpointConnectorName,
		cm.config.GetMongoConnectionString(),
		cm.config.MongoConfig.CheckpointDatabase,
		fmt.Sprintf("%s\\\\.%s", cm.config.MongoConfig.CheckpointDatabase, cm.config.MongoConfig.CheckpointCollection),
	)

	endpoint := fmt.Sprintf("%s/connectors", cm.ConnectionStrings.Debezium)

	for attempt := 1; attempt <= cm.config.DebeziumConfig.MaxRetries; attempt++ {
		log.Debug().Msgf("[%s] Attempting to create Debezium connector (attempt %d/%d)", cm.config.DebeziumConfig.CheckpointConnectorName, attempt, cm.config.DebeziumConfig.MaxRetries)

		// Create HTTP request
		req, err := http.NewRequest("POST", endpoint, strings.NewReader(configJSON))
		if err != nil {
			return fmt.Errorf("[%s] failed to create HTTP request: %w", cm.config.DebeziumConfig.CheckpointConnectorName, err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create HTTP client with timeout
		client := &http.Client{Timeout: 10 * time.Second}

		// Send HTTP request
		res, err := client.Do(req)

		// Handle response
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to create Debezium connector", cm.config.DebeziumConfig.CheckpointConnectorName)
		}

		if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Debug().Msgf("[%s] Debezium connector created successfully on attempt %d", cm.config.DebeziumConfig.CheckpointConnectorName, attempt)

			res.Body.Close()
			return cm.waitForCheckpointConnectorReady(fmt.Sprintf("%s/%s/status", endpoint, cm.config.DebeziumConfig.CheckpointConnectorName))
		}

		if res != nil {
			body := ""
			if b, readErr := io.ReadAll(res.Body); readErr == nil {
				body = string(b)
			}
			res.Body.Close()

			log.Error().Msgf("[%s] Failed to create Debezium connector, status code: %d, response body: %s", cm.config.DebeziumConfig.CheckpointConnectorName, res.StatusCode, body)
		}

		if attempt < cm.config.DebeziumConfig.MaxRetries {
			time.Sleep(cm.config.DebeziumConfig.RetryInterval * time.Duration(attempt))
		}
	}

	return fmt.Errorf("[%s] failed to create Debezium connector after %d attempts", cm.config.DebeziumConfig.CheckpointConnectorName, cm.config.DebeziumConfig.MaxRetries)
}

func (cm *ContainerManager) waitForCheckpointConnectorReady(endpoint string) error {
	deadline := time.Now().Add(cm.config.DebeziumConfig.ReadyTimeout)

	for time.Now().Before(deadline) {
		log.Debug().Msgf("[%s] Checking Debezium connector status", cm.config.DebeziumConfig.CheckpointConnectorName)

		res, err := http.Get(endpoint)

		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to get Debezium connector status", cm.config.DebeziumConfig.CheckpointConnectorName)

			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to read Debezium connector response body", cm.config.DebeziumConfig.CheckpointConnectorName)

			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
		}

		if strings.Contains(string(body), `"state":"RUNNING"`) {
			log.Debug().Msgf("[%s] Debezium connector is ready", cm.config.DebeziumConfig.CheckpointConnectorName)

			time.Sleep(5 * time.Second)
			return nil
		}

		time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
	}

	return nil
}

func (cm *ContainerManager) setupCheckpointToInvocationDebeziumConnector() error {
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

			"transforms": "unwrap,filter,router",

			"transforms.unwrap.type": "io.debezium.connector.mongodb.transforms.ExtractNewDocumentState",

			"transforms.filter.type": "io.debezium.transforms.Filter",
			"transforms.filter.language": "jsr223.groovy",
			"transforms.filter.condition": "value.status == 'SUCCESS'",

			"transforms.router.type": "io.debezium.transforms.ByLogicalTableRouter",
			"transforms.router.topic.regex": "cdc\\.%s\\.%s",
			"transforms.router.topic.replacement": "cdc\\.%s\\.%s"
			}
	}`,
		cm.config.DebeziumConfig.CheckpointToInvocationConnectorName,
		cm.config.GetMongoConnectionString(),
		cm.config.MongoConfig.CheckpointDatabase,
		fmt.Sprintf("%s\\\\.%s", cm.config.MongoConfig.CheckpointDatabase, cm.config.MongoConfig.CheckpointCollection),
		cm.config.MongoConfig.CheckpointDatabase,
		cm.config.MongoConfig.CheckpointCollection,
		cm.config.MongoConfig.InvocationDatabase,
		cm.config.MongoConfig.InvocationCollection,
	)

	endpoint := fmt.Sprintf("%s/connectors", cm.ConnectionStrings.Debezium)

	for attempt := 1; attempt <= cm.config.DebeziumConfig.MaxRetries; attempt++ {
		log.Debug().Msgf("[%s] Attempting to create Debezium connector (attempt %d/%d)", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName, attempt, cm.config.DebeziumConfig.MaxRetries)

		// Create HTTP request
		req, err := http.NewRequest("POST", endpoint, strings.NewReader(configJSON))
		if err != nil {
			return fmt.Errorf("[%s] failed to create HTTP request: %w", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName, err)
		}
		req.Header.Set("Content-Type", "application/json")

		// Create HTTP client with timeout
		client := &http.Client{Timeout: 10 * time.Second}

		// Send HTTP request
		res, err := client.Do(req)

		// Handle response
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to create Debezium connector", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName)
		}

		if err == nil && res.StatusCode >= 200 && res.StatusCode < 300 {
			log.Debug().Msgf("[%s] Debezium connector created successfully on attempt %d", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName, attempt)

			res.Body.Close()
			return cm.waitForCheckpointToInvocationConnectorReady(fmt.Sprintf("%s/%s/status", endpoint, cm.config.DebeziumConfig.CheckpointToInvocationConnectorName))
		}

		if res != nil {
			body := ""
			if b, readErr := io.ReadAll(res.Body); readErr == nil {
				body = string(b)
			}
			res.Body.Close()

			log.Error().Msgf("[%s] Failed to create Debezium connector, status code: %d, response body: %s", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName, res.StatusCode, body)
		}

		if attempt < cm.config.DebeziumConfig.MaxRetries {
			time.Sleep(cm.config.DebeziumConfig.RetryInterval * time.Duration(attempt))
		}
	}

	return fmt.Errorf("[%s] failed to create Debezium connector after %d attempts", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName, cm.config.DebeziumConfig.MaxRetries)
}

func (cm *ContainerManager) waitForCheckpointToInvocationConnectorReady(endpoint string) error {
	deadline := time.Now().Add(cm.config.DebeziumConfig.ReadyTimeout)

	for time.Now().Before(deadline) {
		log.Debug().Msgf("[%s] Checking Debezium connector status", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName)

		res, err := http.Get(endpoint)

		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to get Debezium connector status", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName)

			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			log.Error().Err(err).Msgf("[%s] Failed to read Debezium connector response body", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName)

			time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
		}

		if strings.Contains(string(body), `"state":"RUNNING"`) {
			log.Debug().Msgf("[%s] Debezium connector is ready", cm.config.DebeziumConfig.CheckpointToInvocationConnectorName)

			time.Sleep(5 * time.Second)
			return nil
		}

		time.Sleep(cm.config.DebeziumConfig.ReadyInterval)
	}

	return nil
}

func (cm *ContainerManager) Down() error {
	log.Debug().Msg("Tearing down container manager")

	g := new(errgroup.Group)

	g.Go(func() error {
		if cm.Composes.Kafka != nil {
			if err := cm.Composes.Kafka.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveVolumes(true)); err != nil {
				return fmt.Errorf("failed to tear down Kafka Docker Compose stack: %w", err)
			}
		}

		if cm.Composes.Zookeeper != nil {
			if err := cm.Composes.Zookeeper.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveVolumes(true)); err != nil {
				return fmt.Errorf("failed to tear down Zookeeper Docker Compose stack: %w", err)
			}
		}

		return nil
	})

	g.Go(func() error {
		if cm.Composes.Services != nil {
			if err := cm.Composes.Services.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveVolumes(true)); err != nil {
				return fmt.Errorf("failed to tear down Services Docker Compose stack: %w", err)
			}
		}

		return nil
	})

	g.Go(func() error {
		if cm.Composes.Mongo != nil {
			if err := cm.Composes.Mongo.Down(cm.ctx, compose.RemoveOrphans(true), compose.RemoveVolumes(true)); err != nil {
				return fmt.Errorf("failed to tear down MongoDB Docker Compose stack: %w", err)
			}
		}

		return nil
	})

	if err := g.Wait(); err != nil {
		return fmt.Errorf("failed to wait for container manager down operations: %w", err)
	}

	return nil
}
