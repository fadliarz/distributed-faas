package domain

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type AccumulatorDomainServiceImpl struct{}

func NewAccumulatorDomainService() AccumulatorDomainService {
	return &AccumulatorDomainServiceImpl{}
}

func (s *AccumulatorDomainServiceImpl) AccumulateAmount(charge *Charge, newAmount valueobject.Amount) error {
	charge.AccumulatedAmount = charge.AccumulatedAmount.Add(newAmount)
	return nil
}

func (s *AccumulatorDomainServiceImpl) GenerateNewCharge(userID valueobject.UserID, serviceID valueobject.ServiceID, timestamp valueobject.Timestamp,
	accumulatedAmount valueobject.Amount) (*Charge, error) {

	return &Charge{
		UserID:            userID,
		ServiceID:         serviceID,
		Timestamp:         timestamp,
		AccumulatedAmount: accumulatedAmount,
	}, nil
}

func (s *AccumulatorDomainServiceImpl) ValidateAndInitiateCharge(charge *Charge, chargeID valueobject.ChargeID) error {
	if charge == nil {
		return fmt.Errorf("charge cannot be nil")
	}

	if charge.ChargeID.String() != "" {
		return fmt.Errorf("charge already initiated with ID: %s", charge.ChargeID)
	}

	if chargeID.String() == "" {
		return fmt.Errorf("charge ID cannot be empty")
	}

	charge.ChargeID = chargeID

	if charge.UserID.String() == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	if charge.ServiceID.String() == "" {
		return fmt.Errorf("service ID cannot be empty")
	}

	if charge.Timestamp.Int64() < 0 {
		return fmt.Errorf("timestamp cannot be negative")
	}

	return nil
}
