package main

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/function-service/domain/application-service/service"
	"github.com/fadliarz/distributed-faas/services/function-service/infrastructure/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CreateFunctionIntegrationTestSuite struct {
	suite.Suite
	functionService *service.FunctionApplicationService
	ctx             context.Context
	networkName     string

	networkContainer testcontainers.Network

	// Mongo
	mongoHost      string
	mongoPort      int
	mongoDBName    string
	collectionName string
	mongoUser      string
	mongoPassword  string
	mongoContainer *mongodb.MongoDBContainer
	mongoClient    *mongo.Client

	// Zookeeper and Kafka
	zookeeperContainer testcontainers.Container
	zookeeperHost      string
	zookeeperPort      int

	// Kafka
	kafkaTopic        string
	kafkaHost         string
	kafkaInternalPort int
	kafkaExternalPort int
	kafkaContainer    testcontainers.Container
	kafkaConsumer     *kafka.Consumer
	kafkaAdmin        *kafka.AdminClient

	// Debezium Connnect
	debeziumConnectContainer testcontainers.Container
}

func (suite *CreateFunctionIntegrationTestSuite) SetupSuite() {
}

func (suite *CreateFunctionIntegrationTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.setupEnv()
	suite.setupNetwork()
	suite.setupMongoDB()
	suite.setupZookeeperAndKafka()
	suite.setupDebeziumConnect()
	suite.setupDebeziumConnector()
	suite.setupDependencies()
}

func (suite *CreateFunctionIntegrationTestSuite) IgnoreTestCreateFunctionPersistsToMongoDB() {
	// Arrange
	cmd := &command.CreateFunctionCommand{
		UserID:        uuid.New().String(),
		SourceCodeURL: "https://github.com/user/repo",
	}

	// Act
	functionID, err := suite.functionService.PersistFunction(cmd)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), functionID)

	collection := suite.mongoClient.Database(suite.mongoDBName).Collection(suite.collectionName)

	var result repository.FunctionEntity
	err = collection.FindOne(suite.ctx, bson.M{"function_id": functionID.String()}).Decode(&result)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), functionID.String(), result.FunctionID)
	assert.Equal(suite.T(), cmd.UserID, result.UserID)
	assert.Equal(suite.T(), cmd.SourceCodeURL, result.SourceCodeURL)
}

func (suite *CreateFunctionIntegrationTestSuite) TestCreateFunctionPublishesToKafka() {
	// Arrange
	testID := uuid.New().String()
	cmd := &command.CreateFunctionCommand{
		UserID:        testID,
		SourceCodeURL: fmt.Sprintf("https://github.com/user/repo/%s", testID),
	}

	// Act
	functionID, err := suite.functionService.PersistFunction(cmd)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), functionID)

	// Assert

	suite.T().Logf("Waiting for CDC events on topic: %s", suite.kafkaTopic)
	suite.T().Logf("Looking for testID: %s or functionID: %s", testID, functionID.String())

	time.Sleep(10 * time.Second) // Wait a bit more for Debezium to create its topics

	metadata, err := suite.kafkaAdmin.GetMetadata(nil, false, 5000)
	assert.NoError(suite.T(), err, "Failed to get Kafka metadata")
	found := false
	for topicName := range metadata.Topics {
		if topicName == suite.kafkaTopic {
			found = true
			break
		}
	}
	assert.Equal(suite.T(), true, found, "Expected to find the Kafka topic for CDC events")

	found = false
	for time.Now().Before(time.Now().Add(45*time.Second)) && !found {
		msg, err := suite.kafkaConsumer.ReadMessage(3 * time.Second)
		if err != nil {
			suite.T().Logf("Error reading message: %v", err)
			continue
		}

		if msg != nil {
			msgValue := string(msg.Value)
			if msgValue != "" && (strings.Contains(msgValue, testID) || strings.Contains(msgValue, functionID.String())) {
				suite.T().Logf("Found matching CDC event for function: %s", functionID.String())
				found = true
				break
			} else {
				suite.T().Logf("Message doesn't contain our test data, continuing...")
			}
		}
	}

	assert.True(suite.T(), found, "Expected to find a CDC event for the created function in Kafka")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func TestCreateFunctionIntegrationSuite(t *testing.T) {
	suite.Run(t, new(CreateFunctionIntegrationTestSuite))
}
