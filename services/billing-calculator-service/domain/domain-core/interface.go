package domain

import "github.com/fadliarz/distributed-faas/common/valueobject"

type BillingCalculatorDomainService interface {
	CalculateBilling(billing *Billing, charges []Charge)
	ValidateAndInitiateBilling(billing *Billing, billingID valueobject.BillingID) error
}
