package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (s *MachineApplicationServiceImpl) PersistCheckpoint(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	if cmd.IsRetry == true {
		return s.persistCheckpointIfIsRetryIsTrue(ctx, cmd)
	}

	return s.persistCheckpointIfIsRetryIsFalse(ctx, cmd)

}

func (s *MachineApplicationServiceImpl) persistCheckpointIfIsRetryIsTrue(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	checkpoint := s.mapper.ProcessInvocationCommandToCheckpoint(cmd)

	// Save checkpoint
	err := s.repositoryManager.Checkpoint.UpdateCheckpointTimestampIfStatusIsPendingAndTimestampLessThanThreshold(ctx, checkpoint, 10)
	if err != nil {
		return "", fmt.Errorf("failed to save checkpoint: %w", err)
	}

	return checkpoint.CheckpointID, nil
}

func (s *MachineApplicationServiceImpl) persistCheckpointIfIsRetryIsFalse(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error) {
	checkpoint := s.mapper.ProcessInvocationCommandToCheckpoint(cmd)

	// Validate and initiate checkpoint
	err := s.service.ValidateAndInitiateCheckpoint(checkpoint, domain.NewCheckpointID(primitive.NewObjectID().Hex()))
	if err != nil {
		return "", fmt.Errorf("failed to validate and initiate checkpoint: %w", err)
	}

	// Save checkpoint
	checkpointID, err := s.repositoryManager.Checkpoint.Save(ctx, checkpoint)
	if err != nil {
		return "", fmt.Errorf("failed to save checkpoint: %w", err)
	}

	return checkpointID, nil
}
