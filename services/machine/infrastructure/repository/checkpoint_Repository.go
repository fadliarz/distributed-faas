package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/machine/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
		return "", err
	}

	return domain.NewCheckpointID(checkpointID), nil
}

func (r *CheckpointRepositoryImpl) UpdateCheckpointTimestampIfRetrying(ctx context.Context, checkpoint *domain.Checkpoint, thresholdInSec int64) error {
	entity, err := r.mapper.Entity(checkpoint)
	if err != nil {
		return fmt.Errorf("failed to map checkpoint to entity: %w", err)
	}

	err = r.repo.UpdateCheckpointTimestampIfRetrying(ctx, entity, thresholdInSec)
	if err != nil {
		return err
	}

	return nil
}

func (r *CheckpointRepositoryImpl) UpdateStatusToSuccess(ctx context.Context, checkpointID domain.CheckpointID, outputURL domain.OutputURL) error {
	primitiveCheckpointID, err := primitive.ObjectIDFromHex(checkpointID.String())
	if err != nil {
		return fmt.Errorf("")
	}

	err = r.repo.UpdateStatusToSuccess(ctx, primitiveCheckpointID, outputURL.String())
	if err != nil {
		return fmt.Errorf("failed to update checkpoint status to SUCCESS: %w", err)
	}

	return nil
}
