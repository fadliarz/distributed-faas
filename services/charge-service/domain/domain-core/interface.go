package domain

import (
	"context"
)

type ChargeDomainService interface {
	ValidateCharge(charge *Charge) error
}

type ChargeAggregator interface {
	AddCharge(ctx context.Context, charge *Charge) error
	Start(ctx context.Context)
	Stop()
}

type ChargeProducer interface {
	SendAggregatedCharges(ctx context.Context, aggregates []*ChargeAggregate) error
}
