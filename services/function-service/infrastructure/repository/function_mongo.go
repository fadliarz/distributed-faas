package repository

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FunctionMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewFunctionMongoRepository() *FunctionMongoRepository {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		panic("MONGO_URI environment variable is not set")
	}

	functionMongoDBName := os.Getenv("FUNCTION_MONGO_DB_NAME")
	if functionMongoDBName == "" {
		panic("FUNCTION_MONGO_DB_NAME environment variable is not set")
	}

	functionMongoCollectionName := os.Getenv("FUNCTION_MONGO_COLLECTION_NAME")
	if functionMongoCollectionName == "" {
		panic("FUNCTION_MONGO_COLLECTION_NAME environment variable is not set")
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic("Failed to connect to MongoDB: " + err.Error())
	}

	collection := client.Database(functionMongoDBName).Collection(functionMongoCollectionName)

	return &FunctionMongoRepository{
		client:     client,
		collection: collection,
	}
}

func (r *FunctionMongoRepository) Save(function *FunctionEntity) error {
	_, err := r.collection.InsertOne(context.Background(), function)

	return err
}
