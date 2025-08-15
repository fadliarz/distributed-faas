package application

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
)

type BillingCalculatorDataMapperImpl struct{}

func NewBillingCalculatorDataMapper() BillingCalculatorDataMapper {
	return &BillingCalculatorDataMapperImpl{}
}

func (m *BillingCalculatorDataMapperImpl) BillingCalculationEventToBilling(event *BillingCalculationEvent) *domain.Billing {
	return &domain.Billing{
		UserID:     valueobject.NewUserID(event.UserID),
		LastBilled: valueobject.NewLastBilled(event.LastBilled),
	}
}
