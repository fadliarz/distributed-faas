package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CronMongoRepository struct {
	collection *mongo.Collection
}

func NewCronMongoRepository(collection *mongo.Collection) *CronMongoRepository {
	return &CronMongoRepository{collection: collection}
}

func (r *CronMongoRepository) UpdateLastBilled(ctx context.Context, beforeTimestamp, afterTimestamp int64) error {
	filter := bson.M{
		"last_billed": beforeTimestamp,
	}
	update := bson.M{
		"$set": bson.M{
			"last_billed": afterTimestamp,
		},
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	log.Debug().Msgf("UpdateLastBilled updated %d documents", result.ModifiedCount)

	log.Debug().Msgf("UpdateLastBilled modified %d documents", result.ModifiedCount)

	return nil
}
