package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GenerateBillingIntegrationTest struct {
	IntegrationTestSuite
}

func (suite *GenerateBillingIntegrationTest) TestUserPersisted() {
	// Arrange
	cronEntity := NewCronEntity(primitive.NewObjectID(), time.Date(time.Now().Year(), time.Now().Month()-2, 0, 0, 0, 0, 0, time.UTC).Unix())
	expectedUpdatedLastBilled := time.Date(time.Now().Year(), time.Now().Month()-1, 0, 0, 0, 0, 0, time.UTC).Unix()

	// Act
	err := suite.arrangeHelper.CreateCronInMongoDB(suite.mongoManager.Client, cronEntity)
	require.NoError(suite.T(), err, "Failed to create cron entity")

	// Assert
	suite.assertionHelper.AssertCronUpdatedInMongoDB(suite.ctx, suite.mongoManager.Client, cronEntity.UserID, expectedUpdatedLastBilled)
	suite.assertionHelper.AssertUpdatedCronPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.Cron, cronEntity.UserID, expectedUpdatedLastBilled)
}

func TestGenerateBillingIntegrationTest(t *testing.T) {
	suite.Run(t, new(GenerateBillingIntegrationTest))
}
