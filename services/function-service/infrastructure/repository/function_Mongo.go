package repository

import (
	"context"

	"github.com/fadliarz/distributed-faas/common"
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
		return "", common.MongoWriteErrorHandler(err, nil)
	}

	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (r *FunctionMongoRepository) FindByUserIDAndFunctionID(ctx context.Context, userID string, functionID primitive.ObjectID) (*FunctionEntity, error) {
	filter := bson.M{
		"_id":     functionID,
		"user_id": userID,
	}

	var function FunctionEntity
	err := r.collection.FindOne(context.Background(), filter).Decode(&function)

	if err != nil && err == mongo.ErrNoDocuments {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &function, nil
}

func (r *FunctionMongoRepository) UpdateSourceCodeURLByUserIDAndFunctionID(ctx context.Context, userID string, functionID primitive.ObjectID, sourceCodeURL string) error {
	filter := bson.M{
		"_id":     functionID,
		"user_id": userID,
	}
	update := bson.M{
		"$set": bson.M{
			"source_code_url": sourceCodeURL,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)

	if result != nil && result.MatchedCount == 0 {
		return domain.NewErrFunctionNotFound(err)
	}

	return err
}
