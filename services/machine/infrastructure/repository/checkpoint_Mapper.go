package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CheckpointDataAccessMapperImpl struct{}

func NewCheckpointDataAccessMapperImpl() CheckpointDataAccessMapper {
	return &CheckpointDataAccessMapperImpl{}
}

func (m *CheckpointDataAccessMapperImpl) Entity(checkpoint *domain.Checkpoint) (*CheckpointEntity, error) {
	checkpointID, err := primitive.ObjectIDFromHex(checkpoint.CheckpointID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid checkpoint ID: %v", err)
	}

	return &CheckpointEntity{
		CheckpointID:  checkpointID,
		FunctionID:    checkpoint.FunctionID.String(),
		UserID:        checkpoint.UserID.String(),
		SourceCodeURL: checkpoint.SourceCodeURL.String(),
		Timestamp:     checkpoint.Timestamp.Int64(),
		Status:        checkpoint.Status.String(),
		OutputURL:     checkpoint.OutputURL.String(),
	}, nil
}

func (m *CheckpointDataAccessMapperImpl) Domain(entity *CheckpointEntity) *domain.Checkpoint {
	return &domain.Checkpoint{
		CheckpointID:  domain.NewCheckpointID(entity.CheckpointID.Hex()),
		FunctionID:    domain.NewFunctionID(entity.FunctionID),
		UserID:        domain.NewUserID(entity.UserID),
		SourceCodeURL: domain.NewSourceCodeURL(entity.SourceCodeURL),
		Timestamp:     domain.NewTimestamp(entity.Timestamp),
		Status:        domain.NewStatus(entity.Status),
		OutputURL:     domain.OutputURL(entity.SourceCodeURL),
	}
}
