package repository

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/retry-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/retry-service/domain/domain-core"
)

type CheckpointRepositoryImpl struct {
	repo *CheckpointMongoRepository
}

func NewCheckpointRepository(repo *CheckpointMongoRepository) application.CheckpointRepository {
	return CheckpointRepositoryImpl{
		repo: repo,
	}
}

func (c CheckpointRepositoryImpl) RetryInvocations(ctx context.Context, threshold domain.Threshold) error {
	err := c.repo.RetryInvocations(ctx, threshold.Int64())
	if err != nil {
		return err
	}

	return nil
}
