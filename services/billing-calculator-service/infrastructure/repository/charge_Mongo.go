package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *ChargeMongoRepository) FindChargesByUserIDAndTimeRange(ctx context.Context, userID string, startTime, endTime int64) ([]*ChargeEntity, error) {
	filter := bson.M{
		"user_id": userID,
		"_id": bson.M{
			"$gte": primitive.NewObjectIDFromTimestamp(time.Unix(0, startTime)),
			"$lt":  primitive.NewDateTimeFromTime(time.Time(time.Unix(0, endTime))),
		},
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
