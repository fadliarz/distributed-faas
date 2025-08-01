package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CheckpointMongoRepository struct {
	collection *mongo.Collection
}

func NewCheckpointMongoRepository(collection *mongo.Collection) *CheckpointMongoRepository {
	return &CheckpointMongoRepository{collection: collection}
}

func (r *CheckpointMongoRepository) Save(ctx context.Context, entity *CheckpointEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return "", common.MongoWriteErrorHandler(err, common.NewMongoErrorMapper().WithErrDuplicateKey(domain.ErrCheckpointAlreadyExists))
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *CheckpointMongoRepository) UpdateCheckpointTimestampIfStatusIsPendingAndTimestampLessThanThreshold(
	ctx context.Context, checkpoint *CheckpointEntity, thresholdInSec int64) error {

	result, err := r.collection.UpdateOne(ctx, bson.M{
		"_id":       checkpoint.CheckpointID,
		"status":    "PENDING",
		"timestamp": bson.M{"$lt": checkpoint.Timestamp - thresholdInSec},
	}, bson.M{
		"$set": bson.M{"timestamp": checkpoint.Timestamp},
	})

	if err != nil {
		return fmt.Errorf("failed to update checkpoint timestamp: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewErrCheckpointAlreadyReprocessed(fmt.Errorf("no checkpoint found with ID %s and status PENDING", checkpoint.CheckpointID.Hex()))
	}

	return nil
}

func (r *CheckpointMongoRepository) UpdateStatusToSuccess(ctx context.Context, checkpointID primitive.ObjectID, outputURL string) error {
	update := bson.M{
		"$set": bson.M{
			"status":     "SUCCESS",
			"output_url": outputURL,
		},
	}

	_, err := r.collection.UpdateByID(ctx, checkpointID, update)

	if err != nil {
		return fmt.Errorf("failed to update checkpoint status to SUCCESS: %w", err)
	}

	return nil
}
