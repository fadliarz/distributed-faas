package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fadliarz/services/invocation-service/domain/domain-core/core"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InvocationMongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
	dbName     string
	collName   string
}

func NewInvocationMongoRepository() *InvocationMongoRepository {
	repo := &InvocationMongoRepository{
		dbName:   "invocation_service",
		collName: "invocations",
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

	// ToDo: use connection string from configuration
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return core.NewDatabaseError("failed to connect to MongoDB", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return core.NewDatabaseError("failed to ping MongoDB", err)
	}

	r.client = client
	r.collection = client.Database(r.dbName).Collection(r.collName)
	log.Println("Successfully connected to MongoDB")

	return nil
}

func (r *InvocationMongoRepository) Save(invocation *InvocationEntity) error {
	if invocation == nil {
		return core.NewInternalError("invocation entity cannot be nil", nil)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if r.client == nil || r.collection == nil {
		return core.NewDatabaseError("database connection not initialized", nil)
	}

	_, err := r.collection.InsertOne(ctx, invocation)
	if err != nil {
		return core.NewDatabaseError(
			fmt.Sprintf("failed to save invocation %s", invocation.InvocationID),
			err,
		)
	}

	return nil
}
