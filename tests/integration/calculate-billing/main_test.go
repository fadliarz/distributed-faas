package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CalculateBillingIntegrationTest struct {
	IntegrationTestSuite
}

func (suite *CalculateBillingIntegrationTest) TestBillingCalculated() {
	// Arrange
	userID1 := primitive.NewObjectID().Hex()
	userID2 := primitive.NewObjectID().Hex()
	serviceID1 := primitive.NewObjectID().Hex()
	serviceID2 := primitive.NewObjectID().Hex()
	timestamp := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC).Unix()
	accumulatedAmount := int64(1)

	lastBilled := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC).Unix()

	entities := []*ChargeEntity{
		// UserID 1
		{
			UserID:            userID1,
			ServiceID:         serviceID1,
			Timestamp:         timestamp,
			AccumulatedAmount: accumulatedAmount,
		},
		{
			UserID:            userID1,
			ServiceID:         serviceID2,
			Timestamp:         timestamp,
			AccumulatedAmount: accumulatedAmount,
		},
		// UserID 2
		{
			UserID:            userID2,
			ServiceID:         serviceID1,
			Timestamp:         timestamp,
			AccumulatedAmount: accumulatedAmount,
		},
		{
			UserID:            userID2,
			ServiceID:         serviceID2,
			Timestamp:         timestamp,
			AccumulatedAmount: accumulatedAmount,
		},
	}

	events := []*CronEvent{
		{
			UserID:     userID1,
			LastBilled: lastBilled,
		},
		{
			UserID:     userID2,
			LastBilled: lastBilled,
		},
	}

	// Act
	err := suite.arrangeHelper.CreateCharges(suite.ctx, suite.mongoManager.Client, entities)
	require.NoError(suite.T(), err, "Failed to create charges")
	err = suite.arrangeHelper.CreateCronEvents(suite.ctx, suite.kafkaManager.Producers.Cron, events)
	require.NoError(suite.T(), err, "Failed to create cron events")

	// Assert
	suite.assertionHelper.AssertBillingPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, userID1, accumulatedAmount*2)
	suite.assertionHelper.AssertBillingPersistedInMongoDB(suite.ctx, suite.mongoManager.Client, userID2, accumulatedAmount*2)
}

func TestCalculatedBillingIntegrationTest(t *testing.T) {
	suite.Run(t, new(CalculateBillingIntegrationTest))
}
