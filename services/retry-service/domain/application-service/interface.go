package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/retry-service/domain/domain-core"
)

// Ports

type CheckpointRepository interface {
	RetryInvocations(ctx context.Context, threshold domain.Threshold) error
}
