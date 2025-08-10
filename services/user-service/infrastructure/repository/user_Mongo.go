package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserMongoRepository struct {
	collection *mongo.Collection
}

func NewUserMongoRepository(collection *mongo.Collection) *UserMongoRepository {
	return &UserMongoRepository{
		collection: collection,
	}
}

func (r *UserMongoRepository) Save(ctx context.Context, user *UserEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return "", err
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}
