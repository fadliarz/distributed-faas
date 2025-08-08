package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

func (s *MachineApplicationServiceImpl) PersistCheckpoint(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	if cmd.Status == domain.Pending.String() {
		return s.persistCheckpoint(ctx, cmd)
	}

	if cmd.Status == domain.Retrying.String() {
		return s.updateCheckpointTimestamp(ctx, cmd)
	}

	return "", fmt.Errorf("unsupported status: %s", cmd.Status)
}

func (s *MachineApplicationServiceImpl) updateCheckpointTimestamp(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	checkpoint := s.mapper.ProcessInvocationCommandToCheckpoint(cmd)

	err := s.repositoryManager.Checkpoint.UpdateCheckpointTimestampIfRetrying(ctx, checkpoint, 10)
	if err != nil {
		return "", fmt.Errorf("failed to save checkpoint: %w", err)
	}

	return checkpoint.CheckpointID, nil
}

func (s *MachineApplicationServiceImpl) persistCheckpoint(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	checkpoint := s.mapper.ProcessInvocationCommandToCheckpoint(cmd)

	err := s.service.ValidateAndInitiateCheckpoint(checkpoint)
	if err != nil {
		return "", fmt.Errorf("failed to validate and initiate checkpoint: %w", err)
	}

	checkpointID, err := s.repositoryManager.Checkpoint.Save(ctx, checkpoint)
	if err != nil {
		return "", fmt.Errorf("failed to save checkpoint: %w", err)
	}

	return checkpointID, nil
}
