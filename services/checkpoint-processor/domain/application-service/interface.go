package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/domain/domain-core"
)

// Ports

type InvocationRepository interface {
	UpdateOutputURLAndStatusToSuccessIfNotSet(ctx context.Context, invocationID domain.InvocationID, outputURL string) error
}

type CheckpointConsumer interface {
	common.Consumer[*CheckpointEvent]
}

// Interfaces

type CheckpointProcessorDataMapper interface{}
