package main

import (
	"context"
	"fmt"
	"time"

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

	deadline := time.Now().Add(30 * time.Second)

	for time.Now().Before(deadline) {
		err = mm.Client.Ping(mm.ctx, nil)
		if err != nil {
			continue
		}

		return nil
	}

	return fmt.Errorf("MongoDB connection timed out after 30 seconds")
}
