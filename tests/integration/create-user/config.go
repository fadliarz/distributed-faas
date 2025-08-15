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

	RequestDtos *RequestDtos
}

type ComposePaths struct {
	Common    string
	Mongo     string
	Zookeeper string
	Kafka     string
	Services  string
}

type ContainerNames struct {
	Mongo       string
	Kafka       string
	Debezium    string
	UserService string
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
	UserDatabase   string
	UserCollection string
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
	User string
	Cron string
}

type KafkaConsumerGroups struct {
	User string
	Cron string
}

type DebeziumConfig struct {
	UserConnectorName string
	CronConnectorName string

	MaxRetries    int
	RetryInterval time.Duration
	ReadyTimeout  time.Duration
	ReadyInterval time.Duration
}

type RequestDtos struct {
	CreateUser *CreateUserDto
}

type CreateUserDto struct {
	Password string
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
			Mongo:       "distributed-faas-mongo",
			Kafka:       "distributed-faas-kafka-broker-1",
			Debezium:    "distributed-faas-debezium-connect",
			UserService: "distributed-faas-user-service",
		},
		ComposeConfig: &ComposeConfig{
			ProjectID: "distributed-faas",
			Profile:   "test-create-user",
		},
		MongoConfig: &MongoConfig{
			ReplicaSet:     "rs0",
			Username:       "admin",
			Password:       "password",
			UserDatabase:   "user-db",
			UserCollection: "user",
			CronDatabase:   "cron-db",
			CronCollection: "cron",
		},
		KafkaConfig: &KafkaConfig{
			BootstrapServers: "localhost:19092",
			Topics: &KafkaTopics{
				User: "cdc.user-db.user",
				Cron: "cdc.cron-db.cron",
			},
			ConsumerGroups: &KafkaConsumerGroups{
				User: uuid.NewString(),
				Cron: uuid.NewString(),
			},
			Timeout: 120 * time.Second,
		},
		DebeziumConfig: &DebeziumConfig{
			UserConnectorName: "user-connector",
			CronConnectorName: "cron-connector",
			MaxRetries:        30,
			RetryInterval:     2 * time.Second,
			ReadyTimeout:      6 * time.Minute,
			ReadyInterval:     2 * time.Second,
		},
		RequestDtos: &RequestDtos{
			CreateUser: &CreateUserDto{
				Password: "test-password-123",
			},
		},
	}
}

func (c *TestConfig) GetKafkaConnectionString() string {
	return c.KafkaConfig.BootstrapServers
}

func (c *TestConfig) GetMongoConnectionString() string {
	return "mongodb://" + c.MongoConfig.Username + ":" + c.MongoConfig.Password + "@" + c.ContainerNames.Mongo + ":27017/?replicaSet=" + c.MongoConfig.ReplicaSet + "&directConnection=true"
}
