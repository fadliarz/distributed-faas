package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoManager struct {
	ctx    context.Context
	config *TestConfig

	Client *mongo.Client
}

func NewMongoManager(ctx context.Context, config *TestConfig) *MongoManager {
	return &MongoManager{
		ctx:    ctx,
		config: config,
	}
}

func (mm *MongoManager) SetupClient(mongoURI string) error {
	var err error

	mm.Client, err = mongo.Connect(mm.ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	err = mm.Client.Ping(mm.ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	return nil
}
