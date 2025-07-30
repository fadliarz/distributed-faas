package repository

import (
	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type FunctionEntity struct {
	FunctionID    primitive.ObjectID `bson:"_id,omitempty"`
	UserID        string             `bson:"user_id"`
	SourceCodeURL string             `bson:"source_code_url"`
}

// Mapper

type FunctionDataAccessMapper interface {
	Entity(function *domain.Function) (*FunctionEntity, error)
	Domain(entity *FunctionEntity) *domain.Function
}
