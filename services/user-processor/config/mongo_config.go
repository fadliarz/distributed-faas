package config

import (
	"fmt"
	"os"
)

type CronMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewCronMongoConfig() (*CronMongoConfig, error) {
	mongoURI := os.Getenv("CRON_MONGO_URI")
	if mongoURI == "" {
		return nil, fmt.Errorf("CRON_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("CRON_MONGO_DATABASE")
	if database == "" {
		return nil, fmt.Errorf("CRON_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("CRON_MONGO_COLLECTION")
	if collection == "" {
		return nil, fmt.Errorf("CRON_MONGO_COLLECTION environment variable is not set")
	}

	return &CronMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}
