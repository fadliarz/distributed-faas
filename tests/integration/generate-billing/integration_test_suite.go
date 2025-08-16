package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	ctx              context.Context
	cancel           context.CancelFunc
	config           *TestConfig
	assertionHelper  *AssertionHelper
	arrangeHelper    *ArrangeHelper
	containerManager *ContainerManager
	mongoManager     *MongoManager
	kafkaManager     *KafkaManager
	signalChan       chan os.Signal
	tearDownOnce     bool
}

func (suite *IntegrationTestSuite) SetupSuite() {
	suite.init()

	// Context
	suite.ctx, suite.cancel = context.WithCancel(context.Background())

	// Configuration and Managers
	suite.config = NewDefaultTestConfig()
	suite.assertionHelper = NewAssertionHelper(suite.T(), suite.config)
	suite.arrangeHelper = NewArrangeHelper(suite.T(), suite.config)
	suite.containerManager = NewContainerManager(suite.ctx, suite.config)
	suite.mongoManager = NewMongoManager(suite.ctx, suite.config)
	suite.kafkaManager = NewKafkaManager(suite.ctx, suite.config)

	// Signal handling
	suite.signalChan = make(chan os.Signal, 1)
	signal.Notify(suite.signalChan, syscall.SIGINT, syscall.SIGTERM)
	go suite.handleSignals()

	suite.setupInfrastructure()
}

func (suite *IntegrationTestSuite) init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func (suite *IntegrationTestSuite) TearDownSuite() {
	// Cancel context first to stop the signal handler
	if suite.cancel != nil {
		suite.cancel()
	}

	// Stop signal handling and close the channel safely
	if suite.signalChan != nil {
		signal.Stop(suite.signalChan)
		// Give a moment for the signal handler to exit
		select {
		case <-suite.signalChan:
			// Drain any pending signals
		default:
		}
		close(suite.signalChan)
		suite.signalChan = nil
	}

	// Tear down infrastructure
	if suite.containerManager != nil {
		err := suite.containerManager.Down()
		if err != nil {
			suite.T().Logf("Error tearing down container manager: %v", err)
		}
	}
}

func (suite *IntegrationTestSuite) handleSignals() {
	defer func() {
		// Recover from any panic in signal handling
		if r := recover(); r != nil {
			suite.T().Logf("Signal handler recovered from panic: %v", r)
		}
	}()

	select {
	case sig, ok := <-suite.signalChan:
		if !ok {
			// Channel was closed, exit gracefully
			return
		}
		if sig != nil {
			suite.T().Logf("Received signal: %v, cleaning up...", sig)
			suite.TearDownSuite()
			os.Exit(1)
		}
	case <-suite.ctx.Done():
		return
	}
}

func (suite *IntegrationTestSuite) SetupTest() {
	if suite.ctx == nil {
		suite.ctx = context.Background()
	}
}

func (suite *IntegrationTestSuite) setupInfrastructure() {
	err := suite.containerManager.SetupContainers()
	suite.Require().NoError(err, "Failed to setup infrastructure")

	err = suite.mongoManager.SetupClient(suite.containerManager.ConnectionStrings.Mongo)
	suite.Require().NoError(err, "Failed to setup MongoDB client")

	err = suite.kafkaManager.SetupConsumers(suite.config.GetKafkaConnectionString())
	suite.Require().NoError(err, "Failed to setup Kafka consumers")
}
