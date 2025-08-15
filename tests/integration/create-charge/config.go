package main

import (
	"time"

	"github.com/google/uuid"
)

type TestConfig struct {
	ComposePaths   *ComposePaths
	ContainerNames *ContainerNames

	ComposeConfig *ComposeConfig
	MongoConfig   *MongoConfig
	KafkaConfig   *KafkaConfig

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
	Mongo string
	Kafka string
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
	ChargeDatabase   string
	ChargeCollection string
}

type KafkaConfig struct {
	BootstrapServers string
	Topics           *KafkaTopics
	ConsumerGroups   *KafkaConsumerGroups
	Timeout          time.Duration
}

type KafkaTopics struct {
	Charge string
}

type KafkaConsumerGroups struct {
	Charge string
}

type GrpcEndpoints struct {
	ChargeService string
}

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
			Mongo: "distributed-faas-mongo",
			Kafka: "distributed-faas-kafka-broker-1",
		},
		ComposeConfig: &ComposeConfig{
			ProjectID: "distributed-faas",
			Profile:   "test-create-charge",
		},
		MongoConfig: &MongoConfig{
			ReplicaSet:       "rs0",
			Username:         "admin",
			Password:         "password",
			ChargeDatabase:   "charge-db",
			ChargeCollection: "charge",
		},
		KafkaConfig: &KafkaConfig{
			BootstrapServers: "localhost:19092",
			Topics: &KafkaTopics{
				Charge: "charge",
			},
			ConsumerGroups: &KafkaConsumerGroups{
				Charge: uuid.NewString(),
			},
			Timeout: 120 * time.Second,
		},
		GrpcEndpoints: &GrpcEndpoints{
			ChargeService: "localhost:50059",
		},
	}
}

func (c *TestConfig) GetKafkaConnectionString() string {
	return c.KafkaConfig.BootstrapServers
}

func (c *TestConfig) GetMongoConnectionString() string {
	return "mongodb://" + c.MongoConfig.Username + ":" + c.MongoConfig.Password + "@" + c.ContainerNames.Mongo + ":27017/?replicaSet=" + c.MongoConfig.ReplicaSet + "&directConnection=true"
}
