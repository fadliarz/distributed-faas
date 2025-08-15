package domain

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type BillingCalculatorDomainServiceImpl struct{}

func NewBillingCalculatorDomainService() BillingCalculatorDomainService {
	return &BillingCalculatorDomainServiceImpl{}
}

func (s *BillingCalculatorDomainServiceImpl) CalculateBilling(billing *Billing, charges []Charge) {
	totalAmount := valueobject.NewAmount(0)
	for _, charge := range charges {
		totalAmount = totalAmount.Add(charge.AccumulatedAmount)
	}

	billing.Amount = totalAmount
}

func (s *BillingCalculatorDomainServiceImpl) ValidateAndInitiateBilling(billing *Billing, billingID valueobject.BillingID) error {
	if billing == nil {
		return fmt.Errorf("billing cannot be nil")
	}

	if billing.BillingID.String() != "" {
		return fmt.Errorf("billing ID must be empty for new billing")
	}

	billing.BillingID = billingID

	if billing.UserID.String() == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	if billing.Amount.Int64() < 0 {
		return fmt.Errorf("amount cannot be negative")
	}

	return nil
}
