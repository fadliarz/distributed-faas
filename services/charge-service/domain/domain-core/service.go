package domain

import (
	"strings"
)

// ChargeDomainServiceImpl implements the ChargeDomainService interface
type ChargeDomainServiceImpl struct{}

func NewChargeDomainService() ChargeDomainService {
	return &ChargeDomainServiceImpl{}
}

func (s *ChargeDomainServiceImpl) ValidateCharge(charge *Charge) error {
	if charge == nil {
		return ErrInvalidChargeAmount
	}

	if charge.Amount <= 0 {
		return ErrInvalidChargeAmount
	}

	if strings.TrimSpace(charge.UserID.String()) == "" {
		return ErrInvalidUserID
	}

	if strings.TrimSpace(charge.ServiceID.String()) == "" {
		return ErrInvalidServiceID
	}

	return nil
}
