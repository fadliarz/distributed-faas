package config

import (
	"errors"
	"os"
)

type MachineMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewMachineMongoConfig() (*MachineMongoConfig, error) {
	mongoURI := os.Getenv("MACHINE_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("MACHINE_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("MACHINE_MONGO_DATABASE")
	if database == "" {
		return nil, errors.New("MACHINE_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("MACHINE_MONGO_COLLECTION")
	if collection == "" {
		return nil, errors.New("MACHINE_MONGO_COLLECTION environment variable is not set")
	}

	return &MachineMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}
