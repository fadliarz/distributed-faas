package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/domain-core"
)

// Ports

type ChargeRepository interface {
	UpsertCharges(ctx context.Context, charges []*domain.Charge) error
}

type ChargeConsumer interface {
	PollAndProcessMessages()
}

// Interfaces

type ChargeApplicationService interface {
	ProcessChargeEventBatch(ctx context.Context, events []*ChargeEvent) error
}
