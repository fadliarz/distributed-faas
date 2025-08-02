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
	createFunctionResponse := suite.arrangeHelper.CreateFunction()
	suite.arrangeHelper.UpdateFunctionSourceCodeURL(createFunctionResponse.UserId, createFunctionResponse.FunctionId, "user-id-123/main.js")
	suite.arrangeHelper.RegisterMachine()

	// Act
	createInvocationResponse, err := suite.arrangeHelper.CreateInvocation(createFunctionResponse.UserId, createFunctionResponse.FunctionId)

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
	suite.arrangeHelper.RegisterMachine()

	// Act
	createInvocationResponse, err := suite.arrangeHelper.CreateInvocation(uuid.NewString(), primitive.NewObjectID().Hex())

	// Assert
	suite.assertionHelper.AssertInvocationUnauthorized(createInvocationResponse, err)
}

func TestCreateFunctionCDCIntegrationTest(t *testing.T) {
	suite.Run(t, new(CreateFunctionCDCIntegrationTest))
}
