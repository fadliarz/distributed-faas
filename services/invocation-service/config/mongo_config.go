package config

import (
	"fmt"
	"os"
)

// Invocation

type InvocationMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewInvocationMongoConfig() (*InvocationMongoConfig, error) {
	mongoURI := os.Getenv("INVOCATION_MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("INVOCATION_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("INVOCATION_MONGO_DATABASE")
	if database == "" {
		return nil, fmt.Errorf("INVOCATION_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("INVOCATION_MONGO_COLLECTION")
	if collection == "" {
		return nil, fmt.Errorf("INVOCATION_MONGO_COLLECTION environment variable is not set")
	}

	return &InvocationMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}

// Function Mongo

type FunctionMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewFunctionMongoConfig() (*FunctionMongoConfig, error) {
	mongoURI := os.Getenv("FUNCTION_MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("FUNCTION_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("FUNCTION_MONGO_DATABASE")
	if database == "" {
		return nil, fmt.Errorf("FUNCTION_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("FUNCTION_MONGO_COLLECTION")
	if collection == "" {
		return nil, fmt.Errorf("FUNCTION_MONGO_COLLECTION environment variable is not set")
	}

	return &FunctionMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}
