package application

import (
	"context"
	"time"

	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

// Ports

type CheckpointRepository interface {
	Save(ctx context.Context, checkpoint *domain.Checkpoint) (domain.CheckpointID, error)
	UpdateCheckpointTimestampIfStatusIsPendingAndTimestampLessThanThreshold(
		ctx context.Context, checkpoint *domain.Checkpoint, thresholdInSec int64) error
	UpdateStatusToSuccess(ctx context.Context, checkpointID domain.CheckpointID, outputURL domain.OutputURL) error
}

// Interfaces

type MachineDataMapper interface {
	ProcessInvocationCommandToCheckpoint(cmd *ProcessInvocationCommand) *domain.Checkpoint
}

type MachineApplicationService interface {
	PersistCheckpoint(ctx context.Context, cmd *ProcessInvocationCommand) (domain.CheckpointID, error)
	ExecuteFunction(ctx context.Context, url string, functionID, invocationID string) error
}

// Structs

type LogLine struct {
	Timestamp time.Time
	Content   string
	Source    string
}
