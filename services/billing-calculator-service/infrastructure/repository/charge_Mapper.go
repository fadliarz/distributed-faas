package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChargeDataAccessMapperImpl struct{}

func NewChargeDataAccessMapper() ChargeDataAccessMapper {
	return &ChargeDataAccessMapperImpl{}
}

func (m *ChargeDataAccessMapperImpl) Entity(charge *domain.Charge) (*ChargeEntity, error) {
	if charge == nil {
		return nil, nil
	}

	chargeID, err := primitive.ObjectIDFromHex(charge.ChargeID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid charge ID: %w", err)
	}

	return &ChargeEntity{
		ChargeID:  chargeID,
		UserID:    charge.UserID.String(),
		ServiceID: charge.ServiceID.String(),
		Amount:    charge.Amount.Int64(),
	}, nil
}

func (m *ChargeDataAccessMapperImpl) Domain(entity *ChargeEntity) *domain.Charge {
	if entity == nil {
		return nil
	}

	return domain.NewCharge(
		valueobject.NewChargeID(entity.ChargeID.Hex()),
		valueobject.NewUserID(entity.UserID),
		valueobject.NewServiceID(entity.ServiceID),
		valueobject.NewAmount(entity.Amount),
	)
}
