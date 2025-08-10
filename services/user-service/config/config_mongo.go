package config

import (
	"errors"
	"os"
)

type UserMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewUserMongoConfig() (*UserMongoConfig, error) {
	mongoURI := os.Getenv("USER_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("USER_MONGO_URI environment variable is not set")
	}

	userMongoDBName := os.Getenv("USER_MONGO_DATABASE")
	if userMongoDBName == "" {
		return nil, errors.New("USER_MONGO_DATABASE environment variable is not set")
	}

	userMongoCollectionName := os.Getenv("USER_MONGO_COLLECTION")
	if userMongoCollectionName == "" {
		return nil, errors.New("USER_MONGO_COLLECTION environment variable is not set")
	}

	return &UserMongoConfig{
		MongoURI:   mongoURI,
		Database:   userMongoDBName,
		Collection: userMongoCollectionName,
	}, nil
}
