package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fadliarz/services/invocation-service/domain/domain-core/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvocationMongoRepository struct {
	client     *mongo.Client
	database   string
	collection string
}

func NewInvocationMongoRepository() *InvocationMongoRepository {
	repo := &InvocationMongoRepository{
		database:   os.Getenv("MONGO_DB_DATABASE"),
		collection: os.Getenv("MONGO_DB_INVOCATION_COLLECTION"),
	}

	// ToDo: use connection pooling/singleton
	if err := repo.connect(); err != nil {
		log.Printf("Warning: Failed to connect to MongoDB: %v", err)
	}

	return repo
}

func (r *InvocationMongoRepository) connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_DB_URI")))
	if err != nil {
		return core.NewDatabaseError("failed to connect to MongoDB", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return core.NewDatabaseError("failed to ping MongoDB", err)
	}

	r.client = client

	return nil
}

func (r *InvocationMongoRepository) Save(ctx context.Context, invocation *InvocationEntity) error {
	if invocation == nil {
		return core.NewInternalError("invocation entity cannot be nil", nil)
	}

	if r.client == nil {
		return core.NewDatabaseError("client not initialized", nil)
	}

	_, err := r.client.Database(r.database).Collection(r.collection).InsertOne(ctx, invocation)
	if err != nil {
		return core.NewDatabaseError(
			fmt.Sprintf("failed to save invocation %s", invocation.InvocationID),
			err,
		)
	}

	return nil
}
