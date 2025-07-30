package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/machine/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

type CheckpointRepositoryImpl struct {
	mapper CheckpointDataAccessMapper
	repo   *CheckpointMongoRepository
}

func NewCheckpointRepositoryImpl(mapper CheckpointDataAccessMapper, repo *CheckpointMongoRepository) application.CheckpointRepository {
	return &CheckpointRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *CheckpointRepositoryImpl) Save(ctx context.Context, checkpoint *domain.Checkpoint) (domain.CheckpointID, error) {
	entity, err := r.mapper.Entity(checkpoint)
	if err != nil {
		return "", fmt.Errorf("failed to map checkpoint to entity: %w", err)
	}

	checkpointID, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("failed to save checkpoint: %w", err)
	}

	return domain.NewCheckpointID(checkpointID), nil
}

func (r *CheckpointRepositoryImpl) UpdateCheckpointTimestampIfStatusIsPendingAndTimestampLessThanThreshold(
	ctx context.Context, checkpoint *domain.Checkpoint, thresholdInSec int64) error {

	entity, err := r.mapper.Entity(checkpoint)
	if err != nil {
		return fmt.Errorf("failed to map checkpoint to entity: %w", err)
	}

	err = r.repo.UpdateCheckpointTimestampIfStatusIsPendingAndTimestampLessThanThreshold(ctx, entity, thresholdInSec)
	if err != nil {
		return fmt.Errorf("failed to update checkpoint timestamp: %w", err)
	}

	return nil
}
