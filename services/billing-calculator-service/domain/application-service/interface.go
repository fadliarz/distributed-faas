package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/common"
	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
)

// Ports

type ChargeRepository interface {
	FindChargesByUserIDAndTimeRange(ctx context.Context, userID valueobject.UserID, timestamp valueobject.Timestamp) ([]domain.Charge, error)
}

type BillingRepository interface {
	Save(ctx context.Context, billing *domain.Billing) (valueobject.BillingID, error)
}

type BillingCalculationConsumer interface {
	common.Consumer[*BillingCalculationEvent]
}

// Interfaces

type BillingCalculatorDataMapper interface {
	BillingCalculationEventToBilling(event *BillingCalculationEvent) *domain.Billing
}

type BillingCalculatorApplicationService interface {
	ProcessBillingCalculation(ctx context.Context, event *BillingCalculationEvent) error
}

type BillingCalculationEventHandler interface {
	Handle(ctx context.Context, event *BillingCalculationEvent) error
}
