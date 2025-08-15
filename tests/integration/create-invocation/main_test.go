package main

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateFunctionCDCIntegrationTest struct {
	IntegrationTestSuite
}

func (suite *CreateFunctionCDCIntegrationTest) TestInvocationPersisted() {
	// Arrange
	createFunctionResponse := suite.arrangeHelper.CreateFunction(suite.containerManager.ConnectionStrings.FunctionService)
	suite.arrangeHelper.UpdateFunctionSourceCodeURL(suite.containerManager.ConnectionStrings.FunctionService, createFunctionResponse.UserId, createFunctionResponse.FunctionId, "user-id-123/main.js")
	suite.arrangeHelper.RegisterMachine(suite.containerManager.ConnectionStrings.RegistrarService)

	// Act
	createInvocationResponse, err := suite.arrangeHelper.CreateInvocation(suite.containerManager.ConnectionStrings.InvocationService, createFunctionResponse.UserId, createFunctionResponse.FunctionId)

	// Assert
	suite.assertionHelper.AssertInvocationCreatedSuccessfully(createInvocationResponse, err)
	suite.assertionHelper.AssertInvocationPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
	suite.assertionHelper.AssertInvocationPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.Invocation, createInvocationResponse)
	suite.assertionHelper.AssertCheckpointPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
	suite.assertionHelper.AssertCheckpointPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.Checkpoint, createInvocationResponse)
	suite.assertionHelper.AssertInvocationSuccessInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
}

func (suite *CreateFunctionCDCIntegrationTest) TestInvocation_FunctionNotExists_InvocationUnauthorized() {
	// Arrange
	suite.arrangeHelper.RegisterMachine(suite.containerManager.ConnectionStrings.RegistrarService)

	// Act
	createInvocationResponse, err := suite.arrangeHelper.CreateInvocation(suite.containerManager.ConnectionStrings.InvocationService, uuid.NewString(), primitive.NewObjectID().Hex())

	// Assert
	suite.assertionHelper.AssertInvocationUnauthorized(createInvocationResponse, err)
}

func (suite *CreateFunctionCDCIntegrationTest) TestInvocation_InvocationRetry_InvocationReprocessed() {
	// Arrange
	suite.arrangeHelper.RegisterMachine(suite.containerManager.ConnectionStrings.RegistrarService)

	// Act
	checkpointEntity := suite.arrangeHelper.CreateCheckpointInMongoDB(suite.mongoManager.Client)

	// Assert
	suite.assertionHelper.AssertCheckpointPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, checkpointEntity.CheckpointID.Hex())
}

func TestCreateFunctionCDCIntegrationTest(t *testing.T) {
	suite.Run(t, new(CreateFunctionCDCIntegrationTest))
}
