package application

import (
	"context"
	"fmt"

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
	if err := s.repositoryManager.Checkpoint.RetryInvocations(ctx, threshold); err != nil {
		return fmt.Errorf("failed to retry invocations: %w", err)
	}

	return nil
}
