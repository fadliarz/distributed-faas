package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CronMongoRepository struct {
	collection *mongo.Collection
}

func NewCronMongoRepository(collection *mongo.Collection) *CronMongoRepository {
	return &CronMongoRepository{
		collection: collection,
	}
}

func (r *CronMongoRepository) Save(ctx context.Context, entity *CronEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
