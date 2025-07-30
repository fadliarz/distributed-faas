package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvocationMongoRepository struct {
	collection *mongo.Collection
}

func NewInvocationMongoRepository(collection *mongo.Collection) *InvocationMongoRepository {
	return &InvocationMongoRepository{
		collection: collection,
	}
}

func (r *InvocationMongoRepository) Save(ctx context.Context, invocation *InvocationEntity) (string, error) {
	result, err := r.collection.InsertOne(ctx, invocation)
	if err != nil {
		return "", common.MongoWriteErrorHandler(err, common.NewMongoErrorMapper().WithErrDuplicateKey(domain.NewErrFunctionAlreadyExists(err)))
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *InvocationMongoRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*InvocationEntity, error) {
	var invocation InvocationEntity
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&invocation)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to find invocation by ID: %w", err)
	}

	return &invocation, nil
}
