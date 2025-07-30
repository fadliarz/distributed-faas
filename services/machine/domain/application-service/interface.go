package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

// Ports

type CheckpointRepository interface {
	Save(ctx context.Context, checkpoint *domain.Checkpoint) (domain.CheckpointID, error)
	UpdateCheckpointTimestampIfStatusIsPendingAndTimestampLessThanThreshold(
		ctx context.Context, checkpoint *domain.Checkpoint, thresholdInSec int64) error
}

// Interfaces

type MachineDataMapper interface {
	ProcessInvocationCommandToCheckpoint(cmd *ProcessInvocationCommand) *domain.Checkpoint
}

type MachineApplicationService interface {
	PersistCheckpoint(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error)
}
