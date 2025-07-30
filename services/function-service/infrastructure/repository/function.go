package repository

import (
	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type FunctionEntity struct {
	FunctionID    primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	SourceCodeURL string             `bson:"source_code_url"`
}

// Interfaces

type FunctionDataAccessMapper interface {
	Entity(function *domain.Function) (*FunctionEntity, error)
	Domain(function *FunctionEntity) *domain.Function
}
