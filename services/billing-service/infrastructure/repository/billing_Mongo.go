package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BillingMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewBillingMongoRepository(collection *mongo.Collection) *BillingMongoRepository {
	return &BillingMongoRepository{
		collection: collection,
	}
}

func (r *BillingMongoRepository) FindByUserID(ctx context.Context, userID string) (*BillingEntity, error) {
	filter := bson.M{
		"user_id": userID,
	}

	var billing BillingEntity
	err := r.collection.FindOne(ctx, filter).Decode(&billing)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &billing, nil
}
