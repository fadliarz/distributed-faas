package domain

import (
	"fmt"
	"time"

	"github.com/fadliarz/distributed-faas/common/valueobject"
)

type UserProcessorDomainServiceImpl struct{}

func NewUserProcessorDomainService() UserProcessorDomainService {
	return &UserProcessorDomainServiceImpl{}
}

func (s *UserProcessorDomainServiceImpl) ValidateAndInitiateCron(cron *Cron) error {
	if cron.LastBilling.Int64() != 0 {
		return fmt.Errorf("last billing must be zero for new cron")
	}

	cron.LastBilling = valueobject.NewLastBilled(time.Now().Unix())

	if cron.UserID.String() == "" {
		return fmt.Errorf("invalid user ID")
	}

	return nil
}
