package domain

import (
	"context"
)

// ChargeDomainService defines domain business rules for charges
type ChargeDomainService interface {
	ValidateCharge(charge *Charge) error
}

// ChargeAggregator defines the behavior for charge aggregation
type ChargeAggregator interface {
	AddCharge(ctx context.Context, charge *Charge) error
	Start(ctx context.Context)
	Stop()
}

// ChargeProducer defines the behavior for sending aggregated charges to Kafka
type ChargeProducer interface {
	SendAggregatedCharges(ctx context.Context, aggregates []*ChargeAggregate) error
}
