package repository

import (
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Entity

type CheckpointEntity struct {
	CheckpointID  primitive.ObjectID `bson:"_id,omitempty"`
	InvocationID  string             `bson:"invocation_id"`
	FunctionID    string             `bson:"function_id"`
	SourceCodeURL string             `bson:"source_code_url"`
	Timestamp     int64              `bson:"timestamp"`
	Status        string             `bson:"status"`
	OutputURL     string             `bson:"output_url"`
}

// Mapper

type CheckpointDataAccessMapper interface {
	Entity(checkpoint *domain.Checkpoint) (*CheckpointEntity, error)
	Domain(entity *CheckpointEntity) *domain.Checkpoint
}
