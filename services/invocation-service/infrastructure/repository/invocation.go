package repository

import (
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type InvocationEntity struct {
	InvocationID  primitive.ObjectID `bson:"_id,omitempty"`
	FunctionID    string             `bson:"function_id"`
	UserID        string             `bson:"user_id"`
	SourceCodeURL string             `bson:"source_code_url"`
	OutputURL     string             `bson:"output_url"`
	Status        string             `bson:"status"`
	Timestamp     int64              `bson:"timestamp"`
}

// Mapper

type InvocationDataAccessMapper interface {
	Entity(invocation *domain.Invocation) (*InvocationEntity, error)
	Domain(entity *InvocationEntity) *domain.Invocation
}
