package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type FunctionMongoRepository struct {
	client     *mongo.Client
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
		return "", err
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
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to find function: %w", err)
	}

	return &function, nil
}

func (r *FunctionMongoRepository) UpdateSourceCodeURLByUserIDAndFunctionID(ctx context.Context, userID string, functionID primitive.ObjectID, sourceCodeURL string) error {
	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": functionID, "user_id": userID},
		bson.M{"$set": bson.M{"source_code_url": sourceCodeURL}},
	)

	if err != nil {
		return fmt.Errorf("failed to update function source code URL: %w", err)
	}

	if result.MatchedCount == 0 {
		return domain.NewErrFunctionNotFound(fmt.Errorf("function with ID %s not found", functionID.Hex()))
	}

	return nil
}
