package application

import (
	"time"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
)

type BillingCalculatorDataMapperImpl struct{}

func NewBillingCalculatorDataMapper() BillingCalculatorDataMapper {
	return &BillingCalculatorDataMapperImpl{}
}

func (m *BillingCalculatorDataMapperImpl) BillingCalculationEventToBillingPeriod(event *BillingCalculationEvent) *domain.BillingPeriod {
	lastBilledTime := time.Unix(event.LastBilled, 0)

	return &domain.BillingPeriod{
		StartTimestamp: time.Date(lastBilledTime.Year(), lastBilledTime.Month(), 1, 0, 0, 0, 0, time.UTC).Unix(),
		EndTimestamp:   time.Date(lastBilledTime.Year(), lastBilledTime.Month()+1, 1, 0, 0, 0, 0, time.UTC).Unix(),
		UserID:         valueobject.NewUserID(event.UserID),
	}
}

func (m *BillingCalculatorDataMapperImpl) BillingCalculationEventToBilling(event *BillingCalculationEvent) *domain.Billing {
	return &domain.Billing{
		UserID: valueobject.NewUserID(event.UserID),
	}
}
