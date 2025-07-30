package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrMachineNotFound = errors.New("machine not found")
)

type MachineMongoRepository struct {
	collection *mongo.Collection
}

func NewMachineMongoRepository(collection *mongo.Collection) *MachineMongoRepository {
	return &MachineMongoRepository{
		collection: collection,
	}
}

func (r *MachineMongoRepository) Save(ctx context.Context, entity *MachineEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("failed to insert machine entity: %w", err)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *MachineMongoRepository) UpdateStatus(ctx context.Context, machineID string, address string, status string) error {
	filter := bson.M{"_id": machineID, "address": address}
	update := bson.M{"$set": bson.M{"status": status}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	if result.MatchedCount == 0 {
		return ErrMachineNotFound
	}

	return nil
}
