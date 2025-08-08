package repository

import (
	"context"

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

func (r *InvocationMongoRepository) UpdateOutputURLAndStatusToSuccessIfNotSet(ctx context.Context, invocationID primitive.ObjectID, outputURL string) error {
	filter := bson.M{
		"_id":        invocationID,
		"output_url": "",
	}
	update := bson.M{
		"$set": bson.M{
			"output_url": outputURL,
			"status":     "SUCCESS",
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)

	if result != nil && result.MatchedCount == 0 {
		return nil
	}

	return err
}
