package config

import (
	"fmt"
	"os"
)

type CheckpointMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewCheckpointMongoConfig() (*CheckpointMongoConfig, error) {
	mongoURI := os.Getenv("CHECKPOINT_MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("variable CHECKPOINT_MONGO_URI is not set")
	}

	database := os.Getenv("CHECKPOINT_MONGO_DATABASE")
	if database == "" {
		return nil, fmt.Errorf("variable CHECKPOINT_MONGO_DATABASE is not set")
	}

	collection := os.Getenv("CHECKPOINT_MONGO_COLLECTION")
	if collection == "" {
		return nil, fmt.Errorf("variable CHECKPOINT_MONGO_COLLECTION is not set")
	}

	return &CheckpointMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}
