package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MachineMongoRepository struct {
	collection *mongo.Collection
}

func NewMachineMongoRepository(collection *mongo.Collection) *MachineMongoRepository {
	return &MachineMongoRepository{
		collection: collection,
	}
}

func (r *MachineMongoRepository) FindManyByStatus(ctx context.Context, status string) ([]MachineEntity, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": status})
	if err != nil {
		return nil, fmt.Errorf("failed to find machines by status: %w", err)
	}

	defer cursor.Close(ctx)

	var machines []MachineEntity
	if err = cursor.All(ctx, &machines); err != nil {
		return nil, fmt.Errorf("failed to decode machines: %w", err)
	}

	if machines == nil {
		return []MachineEntity{}, nil
	}

	return machines, nil
}
