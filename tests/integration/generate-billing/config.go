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
}

type ComposePaths struct {
	Common    string
	Mongo     string
	Zookeeper string
	Kafka     string
	Services  string
}

type ContainerNames struct {
	Mongo    string
	Kafka    string
	Debezium string
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
	CronDatabase   string
	CronCollection string
}

type KafkaConfig struct {
	BootstrapServers string
	Topics           *KafkaTopics
	ConsumerGroups   *KafkaConsumerGroups
	Timeout          time.Duration
}

type KafkaTopics struct {
	Cron string
}

type KafkaConsumerGroups struct {
	Cron string
}

type DebeziumConfig struct {
	CronConnectorName string

	MaxRetries    int
	RetryInterval time.Duration
	ReadyTimeout  time.Duration
	ReadyInterval time.Duration
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
			Mongo:    "distributed-faas-mongo",
			Kafka:    "distributed-faas-kafka-broker-1",
			Debezium: "distributed-faas-debezium-connect",
		},
		ComposeConfig: &ComposeConfig{
			ProjectID: "distributed-faas",
			Profile:   "test-generate-billing",
		},
		MongoConfig: &MongoConfig{
			ReplicaSet:     "rs0",
			Username:       "admin",
			Password:       "password",
			CronDatabase:   "cron-db",
			CronCollection: "cron",
		},
		KafkaConfig: &KafkaConfig{
			BootstrapServers: "localhost:19092",
			Topics: &KafkaTopics{
				Cron: "cdc.cron-db.cron",
			},
			ConsumerGroups: &KafkaConsumerGroups{
				Cron: uuid.NewString(),
			},
			Timeout: 120 * time.Second,
		},
		DebeziumConfig: &DebeziumConfig{
			CronConnectorName: "cron-connector",
			MaxRetries:        30,
			RetryInterval:     2 * time.Second,
			ReadyTimeout:      6 * time.Minute,
			ReadyInterval:     2 * time.Second,
		},
	}
}

func (c *TestConfig) GetKafkaConnectionString() string {
	return c.KafkaConfig.BootstrapServers
}

func (c *TestConfig) GetMongoConnectionString() string {
	return "mongodb://" + c.MongoConfig.Username + ":" + c.MongoConfig.Password + "@" + c.ContainerNames.Mongo + ":27017/?replicaSet=" + c.MongoConfig.ReplicaSet + "&directConnection=true"
}
