package main

import (
	"context"
	"testing"

	fs_repository "github.com/fadliarz/distributed-faas/services/function-service/infrastructure/repository"
	is_repository "github.com/fadliarz/distributed-faas/services/invocation-service/infrastructure/repository"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AssertionHelper struct {
	t      *testing.T
	config *TestConfig
}

func NewAssertionHelper(t *testing.T, config *TestConfig) *AssertionHelper {
	return &AssertionHelper{
		t:      t,
		config: config,
	}
}

func (ah *AssertionHelper) AssertFunctionPersistedInFunctionMongoDB(ctx context.Context, client *mongo.Client, functionID string) {
	collection := client.Database(ah.config.MongoConfig.FunctionDatabase).Collection(ah.config.MongoConfig.FunctionCollection)

	objectID, err := primitive.ObjectIDFromHex(functionID)
	require.NoError(ah.t, err, "Failed to convert Function ID to ObjectID")

	var function fs_repository.FunctionEntity
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&function)
	require.NoError(ah.t, err, "Failed to find function in MongoDB")

	require.NotEmpty(ah.t, function.FunctionID, "Function ID should not be empty")
	require.NotEmpty(ah.t, function.UserID, "User ID should not be empty")
	require.Empty(ah.t, function.SourceCodeURL, "")
}

func (ah *AssertionHelper) AssertInvocationPersistedInMongoDB(ctx context.Context, client *mongo.Client, invocationID string) {
	collection := client.Database(ah.config.MongoConfig.InvocationDatabase).Collection(ah.config.MongoConfig.InvocationCollection)

	objectID, err := primitive.ObjectIDFromHex(invocationID)
	require.NoError(ah.t, err, "Failed to convert Invocation ID to ObjectID")

	var invocation is_repository.InvocationEntity
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&invocation)
	require.NoError(ah.t, err, "Failed to find invocation in MongoDB")

	require.NotEmpty(ah.t, invocation.InvocationID, "Invocation ID should not be empty")
	require.NotEmpty(ah.t, invocation.FunctionID, "Function ID should not be empty")
	require.NotEmpty(ah.t, invocation.UserID, "User ID should not be empty")
	require.Empty(ah.t, invocation.SourceCodeURL, "Source code URL should be empty")
	require.Empty(ah.t, invocation.OutputURL, "Output URL should be empty")
	require.False(ah.t, invocation.IsRetry, "Invocation should not be a retry")
	require.Greater(ah.t, invocation.Timestamp, int64(0), "Timestamp should be greater than 0")
}

func contains(s, substr string) bool {
	return len(substr) > 0 && len(s) >= len(substr) &&
		(s == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
