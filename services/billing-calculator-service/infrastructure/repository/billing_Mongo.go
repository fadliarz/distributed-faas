package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BillingMongoRepository struct {
	collection *mongo.Collection
}

func NewBillingMongoRepository(collection *mongo.Collection) *BillingMongoRepository {
	return &BillingMongoRepository{
		collection: collection,
	}
}

func (r *BillingMongoRepository) Save(ctx context.Context, billing *BillingEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, billing)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
