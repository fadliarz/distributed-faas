package domain

import (
	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type AccumulatorDomainService interface {
	AccumulateAmount(charge *Charge, newAmount valueobject.Amount) error
	GenerateNewCharge(userID valueobject.UserID, serviceID valueobject.ServiceID, timestamp valueobject.Timestamp,
		accumulatedAmount valueobject.Amount) (*Charge, error)
	ValidateAndInitiateCharge(charge *Charge, chargeID valueobject.ChargeID) error
}
