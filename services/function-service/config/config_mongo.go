package config

import (
	"errors"
	"os"
)

type FunctionMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewFunctionMongoConfig() (*FunctionMongoConfig, error) {
	mongoURI := os.Getenv("FUNCTION_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("FUNCTION_MONGO_URI environment variable is not set")
	}

	functionMongoDBName := os.Getenv("FUNCTION_MONGO_DATABASE")
	if functionMongoDBName == "" {
		return nil, errors.New("FUNCTION_MONGO_DATABASE environment variable is not set")
	}

	functionMongoCollectionName := os.Getenv("FUNCTION_MONGO_COLLECTION")
	if functionMongoCollectionName == "" {
		return nil, errors.New("FUNCTION_MONGO_COLLECTION environment variable is not set")
	}

	return &FunctionMongoConfig{
		MongoURI:   mongoURI,
		Database:   functionMongoDBName,
		Collection: functionMongoCollectionName,
	}, nil
}
