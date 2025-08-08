package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/retry-service/domain/domain-core"
)

// Constructs

type RetryApplicationService struct {
	repositoryManager *RetryApplicationServiceRepositoryManager
}

type RetryApplicationServiceRepositoryManager struct {
	Checkpoint CheckpointRepository
}

func NewRetryApplicationService(repositoryManager RetryApplicationServiceRepositoryManager) *RetryApplicationService {
	return &RetryApplicationService{
		repositoryManager: &repositoryManager,
	}
}

func NewRetryApplicationServiceRepositoryManager(checkpoint CheckpointRepository) *RetryApplicationServiceRepositoryManager {
	return &RetryApplicationServiceRepositoryManager{
		Checkpoint: checkpoint,
	}
}

// Methods

func (s *RetryApplicationService) RetryInvocations(ctx context.Context, threshold domain.Threshold) error {
	return s.repositoryManager.Checkpoint.RetryInvocations(ctx, threshold)
}
