package application

import (
	"github.com/fadliarz/distributed-faas/services/charge-service/domain/domain-core"
)

// ChargeDataMapperImpl implements the ChargeDataMapper interface
type ChargeDataMapperImpl struct{}

func NewChargeDataMapper() ChargeDataMapper {
	return &ChargeDataMapperImpl{}
}

func (m *ChargeDataMapperImpl) CreateChargeCommandToCharge(cmd *CreateChargeCommand) (*domain.Charge, error) {
	userID := domain.NewUserID(cmd.UserID)
	serviceID := domain.NewServiceID(cmd.ServiceID)
	amount := domain.NewAmount(cmd.Amount)

	return domain.NewCharge(userID, serviceID, amount), nil
}
