package repository

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CheckpointMongoRepository struct {
	collection *mongo.Collection
}

func NewCheckpointMongoRepository(collection *mongo.Collection) *CheckpointMongoRepository {
	return &CheckpointMongoRepository{collection: collection}
}

func (r *CheckpointMongoRepository) RetryInvocations(ctx context.Context, thresholdInSec int64) error {
	filter := bson.M{"status": "PENDING", "timestamp": bson.M{"$lt": time.Now().Unix() - thresholdInSec}}
	update := bson.M{"$set": bson.M{"threshold": time.Now().Unix(), "is_retry": true}}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	log.Debug().Msgf("RetryInvocations updated %d documents", result.ModifiedCount)

	log.Debug().Msgf("RetryInvocations modified %d documents", result.ModifiedCount)

	return nil
}
