package main

import (
	"time"

	"github.com/google/uuid"
)

type TestConfig struct {
	ComposePaths   *ComposePaths
	ContainerNames *ContainerNames

	ComposeConfig  *ComposeConfig
	MongoConfig    *MongoConfig
	KafkaConfig    *KafkaConfig
	DebeziumConfig *DebeziumConfig

	GrpcEndpoints *GrpcEndpoints
}

type ComposePaths struct {
	Common    string
	Mongo     string
	Zookeeper string
	Kafka     string
	Services  string
}

type ContainerNames struct {
	Mongo           string
	Kafka           string
	Debezium        string
	FunctionService string
}

type ComposeConfig struct {
	ProjectID string
	Profile   string
}

type MongoConfig struct {
	// Shared MongoDB configuration
	ReplicaSet string
	Username   string
	Password   string

	// Individual databases and collections
	FunctionDatabase     string
	FunctionCollection   string
	InvocationDatabase   string
	InvocationCollection string
	MachineDatabase      string
	MachineCollection    string
	CheckpointDatabase   string
	CheckpointCollection string
}

type KafkaConfig struct {
	ConsumerGroup string
	AutoCommit    bool

	InvocationTopic string
	CheckpointTopic string
}

type DebeziumConfig struct {
	InvocationConnectorName             string
	CheckpointConnectorName             string
	CheckpointToInvocationConnectorName string

	MaxRetries    int
	RetryInterval time.Duration
	ReadyTimeout  time.Duration
	ReadyInterval time.Duration
}

type GrpcEndpoints struct {
	FunctionService   string
	InvocationService string
}

// Constructor for TestConfig
func NewDefaultTestConfig() *TestConfig {
	return &TestConfig{
		ComposePaths: &ComposePaths{
			Common:    "/home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/common.yml",
			Mongo:     "/home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/mongo.yml",
			Zookeeper: "/home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/zookeeper.yml",
			Kafka:     "/home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/kafka_cluster.yml",
			Services:  "/home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/services.yml",
		},
		ContainerNames: &ContainerNames{
			Mongo:           "distributed-faas-mongo",
			Kafka:           "distributed-faas-kafka-broker-1",
			Debezium:        "distributed-faas-debezium-connect",
			FunctionService: "distributed-faas-function-service",
		},
		ComposeConfig: &ComposeConfig{
			ProjectID: uuid.NewString(),
			Profile:   "test-function-cdc",
		},
		MongoConfig: &MongoConfig{
			// Function MongoDB configuration
			ReplicaSet: "rs0",
			Username:   "admin",
			Password:   "password",

			FunctionDatabase:     "invocation-db",
			FunctionCollection:   "function",
			InvocationDatabase:   "invocation-db",
			InvocationCollection: "invocation",
			MachineDatabase:      "machine-db",
			MachineCollection:    "machine",
			CheckpointDatabase:   "checkpoint-db",
			CheckpointCollection: "checkpoint",
		},
		KafkaConfig: &KafkaConfig{
			ConsumerGroup: "distributed-faas-group",
			AutoCommit:    false,

			InvocationTopic: "cdc.invocation-db.invocation",
			CheckpointTopic: "cdc.checkpoint-db.checkpoint",
		},
		DebeziumConfig: &DebeziumConfig{
			InvocationConnectorName:             "invocation-cdc",
			CheckpointConnectorName:             "checkpoint-cdc",
			CheckpointToInvocationConnectorName: "checkpoint-to-invocation-cdc",

			MaxRetries:    30,              // Number of retries for Debezium connector creation
			RetryInterval: 3 * time.Second, // Interval between retries
			ReadyTimeout:  6 * time.Second,
			ReadyInterval: 2 * time.Second,
		},
		GrpcEndpoints: &GrpcEndpoints{
			FunctionService:   "localhost:50051",
			InvocationService: "localhost:50053",
		},
	}
}

func (config *TestConfig) GetMongoConnectionString() string {
	return "mongodb://" + config.MongoConfig.Username + ":" + config.MongoConfig.Password + "@" + config.ContainerNames.Mongo + ":27017/?replicaSet=" + config.MongoConfig.ReplicaSet + "&directConnection=true"
}

func (config *TestConfig) GetKafkaConnectionString() string {
	return "localhost:19092"
}
