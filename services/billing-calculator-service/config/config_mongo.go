package config

import (
	"errors"
	"os"
)

// Billing

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

	database := os.Getenv("BILLING_MONGO_DATABASE")
	if database == "" {
		return nil, errors.New("BILLING_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("BILLING_MONGO_COLLECTION")
	if collection == "" {
		return nil, errors.New("BILLING_MONGO_COLLECTION environment variable is not set")
	}

	return &BillingMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}

// Charge

type ChargeMongoConfig struct {
	MongoURI   string
	Database   string
	Collection string
}

func NewChargeMongoConfig() (*ChargeMongoConfig, error) {
	mongoURI := os.Getenv("CHARGE_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("CHARGE_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("CHARGE_MONGO_DATABASE")
	if database == "" {
		return nil, errors.New("CHARGE_MONGO_DATABASE environment variable is not set")
	}

	collection := os.Getenv("CHARGE_MONGO_COLLECTION")
	if collection == "" {
		return nil, errors.New("CHARGE_MONGO_COLLECTION environment variable is not set")
	}

	return &ChargeMongoConfig{
		MongoURI:   mongoURI,
		Database:   database,
		Collection: collection,
	}, nil
}
