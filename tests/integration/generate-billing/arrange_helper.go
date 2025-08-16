package main

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo"
)

type ArrangeHelper struct {
	t      *testing.T
	config *TestConfig
}

func NewArrangeHelper(t *testing.T, config *TestConfig) *ArrangeHelper {
	return &ArrangeHelper{
		t:      t,
		config: config,
	}
}

func (h *ArrangeHelper) CreateCronInMongoDB(client *mongo.Client, cronEntity *CronEntity) error {
	collection := client.Database(h.config.MongoConfig.CronDatabase).Collection(h.config.MongoConfig.CronCollection)

	_, err := collection.InsertOne(h.t.Context(), cronEntity)

	return err
}
