package config

import (
	"errors"
	"os"
)

type BillingMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewBillingMongoConfig() (*BillingMongoConfig, error) {
	mongoURI := os.Getenv("BILLING_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("BILLING_MONGO_URI environment variable is not set")
	}

	billingMongoDBName := os.Getenv("BILLING_MONGO_DATABASE")
	if billingMongoDBName == "" {
		return nil, errors.New("BILLING_MONGO_DATABASE environment variable is not set")
	}

	billingMongoCollectionName := os.Getenv("BILLING_MONGO_COLLECTION")
	if billingMongoCollectionName == "" {
		return nil, errors.New("BILLING_MONGO_COLLECTION environment variable is not set")
	}

	return &BillingMongoConfig{
		MongoURI:   mongoURI,
		Database:   billingMongoDBName,
		Collection: billingMongoCollectionName,
	}, nil
}
