package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvocationMongoRepository struct {
	collection *mongo.Collection
}

func NewInvocationMongoRepository(collection *mongo.Collection) *InvocationMongoRepository {
	return &InvocationMongoRepository{
		collection: collection,
	}
}

func (r *InvocationMongoRepository) UpdateOutputURLIfNotSet(ctx context.Context, invocationID primitive.ObjectID, outputURL string) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":        invocationID,
		"output_url": "",
	}, bson.M{
		"$set": bson.M{"output_url": outputURL},
	})

	if err != nil {
		return fmt.Errorf("failed to update invocation: %w", err)
	}

	if result.MatchedCount == 0 {
		return nil
	}

	return nil
}
