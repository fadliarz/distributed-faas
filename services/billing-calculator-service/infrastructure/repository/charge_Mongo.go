package repository

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ChargeMongoRepository struct {
	collection *mongo.Collection
}

func NewChargeMongoRepository(collection *mongo.Collection) *ChargeMongoRepository {
	return &ChargeMongoRepository{
		collection: collection,
	}
}

func (r *ChargeMongoRepository) FindChargesByUserIDAndTimeRange(ctx context.Context, userID string, timestamp int64) ([]*ChargeEntity, error) {
	log.Debug().
		Str("user_id", userID).
		Int64("timestamp", timestamp).
		Msg("Finding charges")

	filter := bson.M{
		"user_id":   userID,
		"timestamp": timestamp,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var entities []*ChargeEntity
	if err = cursor.All(ctx, &entities); err != nil {
		return nil, err
	}

	if entities == nil {
		return []*ChargeEntity{}, nil
	}

	return entities, nil
}
