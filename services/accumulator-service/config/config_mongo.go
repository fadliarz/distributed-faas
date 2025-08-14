package config

import (
	"errors"
	"os"
)

type AccumulatorMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewAccumulatorMongoConfig() (*AccumulatorMongoConfig, error) {
	mongoURI := os.Getenv("ACCUMULATOR_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("ACCUMULATOR_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("ACCUMULATOR_MONGO_DATABASE")
	if database == "" {
		return nil, errors.New("ACCUMULATOR_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("ACCUMULATOR_MONGO_COLLECTION")
	if collection == "" {
		return nil, errors.New("ACCUMULATOR_MONGO_COLLECTION environment variable is not set")
	}

	return &AccumulatorMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}
