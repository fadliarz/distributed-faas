package config

import (
	"errors"
	"os"
)

type BillingCalculatorMongoConfig struct {
	MongoURI          string
	Database          string
	ChargeCollection  string
	BillingCollection string
}

func NewBillingCalculatorMongoConfig() (*BillingCalculatorMongoConfig, error) {
	mongoURI := os.Getenv("BILLING_CALCULATOR_MONGO_URI")
	if mongoURI == "" {
		return nil, errors.New("BILLING_CALCULATOR_MONGO_URI environment variable is not set")
	}

	database := os.Getenv("BILLING_CALCULATOR_MONGO_DATABASE")
	if database == "" {
		return nil, errors.New("BILLING_CALCULATOR_MONGO_DATABASE environment variable is not set")
	}

	chargeCollection := os.Getenv("BILLING_CALCULATOR_CHARGE_COLLECTION")
	if chargeCollection == "" {
		return nil, errors.New("BILLING_CALCULATOR_CHARGE_COLLECTION environment variable is not set")
	}

	billingCollection := os.Getenv("BILLING_CALCULATOR_BILLING_COLLECTION")
	if billingCollection == "" {
		return nil, errors.New("BILLING_CALCULATOR_BILLING_COLLECTION environment variable is not set")
	}

	return &BillingCalculatorMongoConfig{
		MongoURI:          mongoURI,
		Database:          database,
		ChargeCollection:  chargeCollection,
		BillingCollection: billingCollection,
	}, nil
}
