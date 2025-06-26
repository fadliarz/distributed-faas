package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/docker/go-connections/nat"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/service"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (suite *CreateFunctionIntegrationTestSuite) setupDependencies() {
	suite.functionService = service.NewFunctionApplicationService()
}

func (suite *CreateFunctionIntegrationTestSuite) setupEnv() {
	suite.networkName = "testcontainers-network"

	suite.mongoHost = "mongo-test"
	suite.mongoPort = 27017
	suite.mongoDBName = "faas"
	suite.collectionName = "functions"
	suite.mongoUser = "admin"
	suite.mongoPassword = "password"

	suite.zookeeperHost = "zookeeper-test"
	suite.zookeeperPort = 2181

	suite.kafkaHost = "kafka-test"
	suite.kafkaInternalPort = 9092
	suite.kafkaExternalPort = 29092

	suite.kafkaTopic = fmt.Sprintf("test.%s.%s", suite.mongoDBName, suite.collectionName)

	os.Setenv("FUNCTION_MONGO_DB_NAME", suite.mongoDBName)
	os.Setenv("FUNCTION_MONGO_COLLECTION_NAME", suite.collectionName)
}

func (suite *CreateFunctionIntegrationTestSuite) setupNetwork() {
	networkContainer, err := testcontainers.GenericNetwork(suite.ctx, testcontainers.GenericNetworkRequest{
		NetworkRequest: testcontainers.NetworkRequest{
			Name:           suite.networkName,
			CheckDuplicate: false,
		},
	})
	if err != nil {
		suite.T().Fatalf("Failed to create network: %v", err)
	}
	suite.networkContainer = networkContainer
}

func (suite *CreateFunctionIntegrationTestSuite) setupMongoDB() {
	var err error

	suite.mongoContainer, err = mongodb.Run(
		suite.ctx, "mongo:6",
		mongodb.WithReplicaSet("rs0"),
		mongodb.WithUsername(suite.mongoUser),
		mongodb.WithPassword(suite.mongoPassword),
		testcontainers.WithExposedPorts(fmt.Sprintf("%d/tcp", suite.mongoPort)),
		network.WithNetwork([]string{suite.networkName}, suite.networkContainer.(*testcontainers.DockerNetwork)),
		testcontainers.WithExposedPorts(fmt.Sprintf("%d/tcp", suite.mongoPort)),
	)
	assert.NoError(suite.T(), err, "Failed to start MongoDB container")

	containerInspect, err := suite.mongoContainer.Inspect(suite.ctx)
	assert.NoError(suite.T(), err, "Failed to inspect mongo container")
	suite.mongoHost = containerInspect.Name[1:] // Remove leading slash

	time.Sleep(2 * time.Second)

	mongoURI, err := suite.mongoContainer.ConnectionString(suite.ctx)
	assert.NoError(suite.T(), err, "Failed to get MongoDB connection string")
	mongoURI = strings.Replace(mongoURI, "localhost", "127.0.0.1", 1)
	mongoURI = fmt.Sprintf("%s&directConnection=true", mongoURI) // Ensure replica set is specified

	os.Setenv("MONGO_URI", mongoURI)

	suite.waitForMongoReplicaSet(mongoURI)

	suite.mongoClient, err = mongo.Connect(suite.ctx, options.Client().ApplyURI(mongoURI))
	assert.NoError(suite.T(), err, "Failed to connect to MongoDB")

	err = suite.mongoClient.Database(suite.mongoDBName).Collection(suite.collectionName).Drop(suite.ctx)
	assert.NoError(suite.T(), err, "Failed to drop collection before test")
}

func (suite *CreateFunctionIntegrationTestSuite) waitForMongoReplicaSet(uri string) {
	maxRetries := 10
	retryInterval := 2 * time.Second

	log.Println("Waiting for MongoDB replica set to be ready...")
	for i := 0; i < maxRetries; i++ {
		client, err := mongo.Connect(suite.ctx, options.Client().ApplyURI(uri).SetServerSelectionTimeout(5*time.Second))
		if err != nil {
			suite.T().Logf("Failed to connect to MongoDB (attempt %d/%d): %v", i+1, maxRetries, err)
			time.Sleep(retryInterval)
			continue
		}

		err = client.Ping(suite.ctx, nil)
		client.Disconnect(suite.ctx)
		if err == nil {
			log.Println("MongoDB replica set is ready!")
			return
		}

		suite.T().Logf("MongoDB replica set not ready yet (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(retryInterval)
	}

	suite.T().Fatalf("MongoDB replica set failed to initialize after %d attempts", maxRetries)
}

func (suite *CreateFunctionIntegrationTestSuite) setupZookeeperAndKafka() {
	var err error

	suite.zookeeperContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "bitnami/zookeeper:3.8",
			ExposedPorts: []string{fmt.Sprintf("%d/tcp", suite.zookeeperPort)},
			Env: map[string]string{
				"ALLOW_ANONYMOUS_LOGIN": "yes",
			},
			WaitingFor: wait.ForLog("binding to port").WithStartupTimeout(30 * time.Second),
			Networks:   []string{suite.networkName},
			Name:       suite.zookeeperHost,
		},
		Started: true,
	})
	assert.NoError(suite.T(), err, "Failed to start Zookeeper container")

	suite.kafkaContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image: "bitnami/kafka:3.4",
			ExposedPorts: []string{
				"9092/tcp",
				"29092:29092/tcp",
			},
			Env: map[string]string{
				"KAFKA_CFG_ZOOKEEPER_CONNECT":                fmt.Sprintf("%s:%d", suite.zookeeperHost, suite.zookeeperPort),
				"KAFKA_CFG_LISTENERS":                        "PLAINTEXT://:9092,EXTERNAL://:29092",
				"KAFKA_CFG_ADVERTISED_LISTENERS":             fmt.Sprintf("PLAINTEXT://%s:9092,EXTERNAL://%s:29092", suite.kafkaHost, "localhost"),
				"KAFKA_CFG_LISTENER_SECURITY_PROTOCOL_MAP":   "PLAINTEXT:PLAINTEXT,EXTERNAL:PLAINTEXT",
				"KAFKA_CFG_INTER_BROKER_LISTENER_NAME":       "PLAINTEXT",
				"ALLOW_PLAINTEXT_LISTENER":                   "yes",
				"KAFKA_CFG_AUTO_CREATE_TOPICS_ENABLE":        "true",
				"KAFKA_CFG_OFFSETS_TOPIC_REPLICATION_FACTOR": "1",
			},
			WaitingFor: wait.ForLog("started (kafka.server.KafkaServer)").WithStartupTimeout(60 * time.Second),
			Networks:   []string{suite.networkName},
			Name:       suite.kafkaHost,
		},
		Started: true,
	})
	assert.NoError(suite.T(), err, "Failed to start Kafka container")

	port, err := suite.kafkaContainer.MappedPort(suite.ctx, nat.Port(strconv.Itoa(suite.kafkaExternalPort)))
	assert.NoError(suite.T(), err, "Failed to get mapped port for Kafka container")
	host, err := suite.kafkaContainer.Host(suite.ctx)
	assert.NoError(suite.T(), err, "Failed to get host for Kafka container")

	suite.kafkaConsumer, err = kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s:%d", host, port.Int()),
		"group.id":           "test-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	assert.NoError(suite.T(), err, "Failed to create Kafka consumer")

	suite.kafkaAdmin, err = kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("%s:%d", host, port.Int()),
	})
	assert.NoError(suite.T(), err, "Failed to create Kafka admin client")

	metadata, err := suite.kafkaAdmin.GetMetadata(&suite.kafkaTopic, false, 300000)
	if err != nil {
		suite.T().Logf("Warning: Failed to get metadata for topic %s: %v", suite.kafkaTopic, err)
	} else if metadata != nil {
		suite.T().Logf("Topic %s exists with %d partition(s)",
			suite.kafkaTopic, len(metadata.Topics[suite.kafkaTopic].Partitions))
	}

	err = suite.kafkaConsumer.SubscribeTopics([]string{suite.kafkaTopic}, nil)
	assert.NoError(suite.T(), err, "Failed to subscribe to Kafka topic")
}

func (suite *CreateFunctionIntegrationTestSuite) setupDebeziumConnect() {
	var err error
	suite.debeziumConnectContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "debezium/connect:2.3",
			ExposedPorts: []string{"8083/tcp"},
			Env: map[string]string{
				"BOOTSTRAP_SERVERS":                 fmt.Sprintf("%s:%d", suite.kafkaHost, 9092),
				"GROUP_ID":                          "1",
				"CONFIG_STORAGE_TOPIC":              "connect_configs",
				"OFFSET_STORAGE_TOPIC":              "connect_offsets",
				"STATUS_STORAGE_TOPIC":              "connect_statuses",
				"KEY_CONVERTER":                     "org.apache.kafka.connect.json.JsonConverter",
				"VALUE_CONVERTER":                   "org.apache.kafka.connect.json.JsonConverter",
				"CONNECT_REST_ADVERTISED_HOST_NAME": "connect",
			},
			Networks: []string{suite.networkName},
		},
		Started: true,
	})
	assert.NoError(suite.T(), err, "Failed to start Debezium Connect container")
}
