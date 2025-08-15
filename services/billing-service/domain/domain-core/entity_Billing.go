package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type Billing struct {
	BillingID  valueobject.BillingID
	UserID     valueobject.UserID
	LastBilled valueobject.LastBilled
	Amount     valueobject.Amount
}
