package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChargeMongoRepository struct {
	collection *mongo.Collection
}

func NewChargeMongoRepository(collection *mongo.Collection) *ChargeMongoRepository {
	return &ChargeMongoRepository{
		collection: collection,
	}
}

func (r *ChargeMongoRepository) UpsertCharges(ctx context.Context, charges []*ChargeEntity) error {
	if len(charges) == 0 {
		return nil
	}

	var writeModels []mongo.WriteModel

	for _, charge := range charges {
		filter := bson.M{
			"user_id":    charge.UserID,
			"service_id": charge.ServiceID,
			"timestamp":  charge.Timestamp,
		}
		update := bson.M{"$inc": bson.M{"accumulated_amount": charge.AccumulatedAmount}}

		writeModels = append(writeModels, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	bulkWriteOptions := options.BulkWrite().SetOrdered(false) // Set to false for better performance

	_, err := r.collection.BulkWrite(ctx, writeModels, bulkWriteOptions)

	return err
}
