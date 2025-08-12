package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/charge-service/domain/domain-core"
)

// Ports

type ChargeProducer interface {
	SendAggregatedCharges(ctx context.Context, aggregates []*domain.ChargeAggregate) error
}

type ChargeAggregator interface {
	AddCharge(ctx context.Context, charge *domain.Charge) error
}

// Interfaces

type ChargeDataMapper interface {
	CreateChargeCommandToCharge(cmd *CreateChargeCommand) (*domain.Charge, error)
}

type ChargeApplicationService interface {
	ProcessCharge(ctx context.Context, command *CreateChargeCommand) error
}
