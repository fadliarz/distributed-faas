package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
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
	createInvocationResponse := suite.arrangeHelper.CreateInvocation(createFunctionResponse.UserId, createFunctionResponse.FunctionId)

	// Assert
	suite.assertionHelper.AssertInvocationPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
	suite.assertionHelper.AssertInvocationPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.Invocation, createInvocationResponse)
	suite.assertionHelper.AssertCheckpointPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
	suite.assertionHelper.AssertInvocationSuccessInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
}

func TestCreateFunctionCDCIntegrationTest(t *testing.T) {
	suite.Run(t, new(CreateFunctionCDCIntegrationTest))
}
