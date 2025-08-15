package application

import (
	"time"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/charge-service/domain/domain-core"
)

type ChargeDataMapperImpl struct{}

func NewChargeDataMapper() ChargeDataMapper {
	return &ChargeDataMapperImpl{}
}

func (m *ChargeDataMapperImpl) CreateChargeCommandToCharge(cmd *CreateChargeCommand) (*domain.Charge, error) {
	userID := valueobject.NewUserID(cmd.UserID)
	serviceID := valueobject.NewServiceID(cmd.ServiceID)
	amount := valueobject.NewAmount(cmd.Amount)

	return &domain.Charge{
		UserID:    userID,
		ServiceID: serviceID,
		Amount:    amount,
		Timestamp: valueobject.NewTimestamp(time.Now().Unix()),
	}, nil
}
