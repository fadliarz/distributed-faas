package main

import (
	"testing"

	charge_service_v1 "github.com/fadliarz/distributed-faas/services/charge-service/gen/go/charge-service/v1"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserCDCIntegrationTest struct {
	IntegrationTestSuite
}

func (suite *CreateUserCDCIntegrationTest) TestUserPersisted() {
	// Arrange
	userID1 := primitive.NewObjectID().Hex()
	userID2 := primitive.NewObjectID().Hex()
	serviceID1 := primitive.NewObjectID().Hex()
	serviceID2 := primitive.NewObjectID().Hex()
	amount := int64(1)

	requests := []*charge_service_v1.CreateChargeRequest{
		// UserID 1
		{
			UserId:    userID1,
			ServiceId: serviceID1,
			Amount:    amount,
		},
		{
			UserId:    userID1,
			ServiceId: serviceID1,
			Amount:    amount,
		},
		{
			UserId:    userID1,
			ServiceId: serviceID2,
			Amount:    amount,
		},
		{
			UserId:    userID1,
			ServiceId: serviceID2,
			Amount:    amount,
		},
		// UserID 2
		{
			UserId:    userID2,
			ServiceId: serviceID1,
			Amount:    amount,
		},
		{
			UserId:    userID2,
			ServiceId: serviceID1,
			Amount:    amount,
		},
		{
			UserId:    userID2,
			ServiceId: serviceID2,
			Amount:    amount,
		},
		{
			UserId:    userID2,
			ServiceId: serviceID2,
			Amount:    amount,
		},
	}

	// Act
	err := suite.arrangeHelper.CreateCharges(suite.containerManager.ConnectionStrings.ChargeService, requests)
	require.NoError(suite.T(), err, "Failed to create charges")

	// Assert
	suite.assertionHelper.AssertChargesPersistedInKafka(suite.ctx, suite.kafkaManager.Consumers.Charge, userID1, userID2, serviceID1, serviceID2, 2*amount)
	suite.assertionHelper.AssertChargePersistedInMongoDB(suite.ctx, suite.mongoManager.Client, userID1, serviceID1, 2*amount)
	suite.assertionHelper.AssertChargePersistedInMongoDB(suite.ctx, suite.mongoManager.Client, userID2, serviceID2, 2*amount)
	suite.assertionHelper.AssertChargePersistedInMongoDB(suite.ctx, suite.mongoManager.Client, userID2, serviceID1, 2*amount)
	suite.assertionHelper.AssertChargePersistedInMongoDB(suite.ctx, suite.mongoManager.Client, userID2, serviceID2, 2*amount)
}

func TestCreateUserCDCIntegrationTest(t *testing.T) {
	suite.Run(t, new(CreateUserCDCIntegrationTest))
}
