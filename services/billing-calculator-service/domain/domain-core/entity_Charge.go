package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

// Charge represents a charge record for billing calculation
type Charge struct {
	ChargeID          valueobject.ChargeID
	UserID            valueobject.UserID
	ServiceID         valueobject.ServiceID
	AccumulatedAmount valueobject.Amount
}

func NewCharge(chargeID valueobject.ChargeID, userID valueobject.UserID, serviceID valueobject.ServiceID, accumulatedAmount valueobject.Amount) *Charge {
	return &Charge{
		ChargeID:          chargeID,
		UserID:            userID,
		ServiceID:         serviceID,
		AccumulatedAmount: accumulatedAmount,
	}
}
