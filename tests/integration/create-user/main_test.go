package main

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CreateUserCDCIntegrationTest struct {
	IntegrationTestSuite
}

func (suite *CreateUserCDCIntegrationTest) TestUserPersisted() {
	// Arrange
	password := "test-password-123"

	// Act
	createUserResponse, err := suite.arrangeHelper.CreateUser(password)

	// Assert
	suite.assertionHelper.AssertUserCreatedSuccessfully(createUserResponse, err)
	suite.assertionHelper.AssertUserPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createUserResponse.UserId, password)
	suite.assertionHelper.AssertUserPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.User, createUserResponse)
	suite.assertionHelper.AssertCronJobPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, createUserResponse.UserId)
}

func TestCreateUserCDCIntegrationTest(t *testing.T) {
	suite.Run(t, new(CreateUserCDCIntegrationTest))
}
