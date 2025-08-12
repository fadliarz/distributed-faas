package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

// Charge represents a charge record for billing calculation
type Charge struct {
	ChargeID  valueobject.ChargeID
	UserID    valueobject.UserID
	ServiceID valueobject.ServiceID
	Amount    valueobject.Amount
}

func NewCharge(chargeID valueobject.ChargeID, userID valueobject.UserID, serviceID valueobject.ServiceID, amount valueobject.Amount) *Charge {
	return &Charge{
		ChargeID:  chargeID,
		UserID:    userID,
		ServiceID: serviceID,
		Amount:    amount,
	}
}
