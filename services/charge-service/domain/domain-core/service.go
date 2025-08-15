package domain

import (
	"fmt"
)

type ChargeDomainServiceImpl struct{}

func NewChargeDomainService() ChargeDomainService {
	return &ChargeDomainServiceImpl{}
}

func (s *ChargeDomainServiceImpl) ValidateCharge(charge *Charge) error {
	if charge == nil {
		return fmt.Errorf("charge cannot be nil")
	}

	if charge.Amount <= 0 {
		return fmt.Errorf("charge amount must be greater than zero")
	}

	if charge.UserID.String() == "" {
		return fmt.Errorf("user ID cannot be empty")
	}

	if charge.ServiceID.String() == "" {
		return fmt.Errorf("service ID cannot be empty")
	}

	return nil
}
