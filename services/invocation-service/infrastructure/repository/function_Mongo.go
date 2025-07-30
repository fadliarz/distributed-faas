package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FunctionMongoRepository struct {
	collection *mongo.Collection
}

func NewFunctionMongoRepository(collection *mongo.Collection) *FunctionMongoRepository {
	return &FunctionMongoRepository{
		collection: collection,
	}
}

func (r *FunctionMongoRepository) Save(ctx context.Context, function *FunctionEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, function)
	if err != nil {
		return "", common.MongoWriteErrorHandler(err, nil)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *FunctionMongoRepository) FindByUserIDAndFunctionID(ctx context.Context, userID string, functionID primitive.ObjectID) (*FunctionEntity, error) {
	filter := map[string]interface{}{
		"_id":     functionID,
		"user_id": userID,
	}

	var function FunctionEntity
	err := r.collection.FindOne(context.Background(), filter).Decode(&function)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, domain.NewErrFunctionNotFound(err)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find function: %w", err)
	}

	return &function, nil
}
