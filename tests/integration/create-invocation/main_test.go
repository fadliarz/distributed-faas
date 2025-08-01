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

	// Act
	createInvocationResponse := suite.arrangeHelper.CreateInvocation(createFunctionResponse.UserId, createFunctionResponse.FunctionId)

	// Assert
	suite.assertionHelper.AssertInvocationPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createInvocationResponse.InvocationId)
	suite.assertionHelper.AssertInvocationPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.Invocation, createInvocationResponse)
}

func TestCreateFunctionCDCIntegrationTest(t *testing.T) {
	suite.Run(t, new(CreateFunctionCDCIntegrationTest))
}
