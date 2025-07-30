package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/domain/domain-core"
)

// Constructors

type CheckpointEventHandler struct {
	mapper            CheckpointProcessorDataMapper
	repositoryManager *CheckpointEventHandlerRepositoryManager
}

type CheckpointEventHandlerRepositoryManager struct {
	Invocation InvocationRepository
}

func NewCheckpointEventHandler(mapper CheckpointProcessorDataMapper, repositoryManager *CheckpointEventHandlerRepositoryManager) *CheckpointEventHandler {
	return &CheckpointEventHandler{
		mapper:            mapper,
		repositoryManager: repositoryManager,
	}
}

func NewCheckpointEventHandlerRepositoryManager(invocation InvocationRepository) *CheckpointEventHandlerRepositoryManager {
	return &CheckpointEventHandlerRepositoryManager{
		Invocation: invocation,
	}
}

// Methods

func (eh *CheckpointEventHandler) HandleCheckpointEvent(ctx context.Context, event *CheckpointEvent) error {
	err := eh.repositoryManager.Invocation.UpdateOutputURLIfNotSet(ctx, domain.InvocationID(event.CheckpointID), event.OutputURL)
	if err != nil {
		return fmt.Errorf("failed to update output URL: %w", err)
	}

	return nil
}
