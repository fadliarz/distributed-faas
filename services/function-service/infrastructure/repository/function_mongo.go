package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

type FunctionMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewFunctionMongoRepository() *FunctionMongoRepository {
	return &FunctionMongoRepository{}
}

func (r *FunctionMongoRepository) Save(function *FunctionEntity) error {
	defaultErr := errors.New("")

	_, err := r.collection.InsertOne(context.Background(), function)
	if err != nil {
		return defaultErr
	}

	return nil
}
