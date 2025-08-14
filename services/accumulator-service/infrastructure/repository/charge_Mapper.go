package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChargeDataAccessMapperImpl struct{}

func NewChargeDataAccessMapper() ChargeDataAccessMapper {
	return &ChargeDataAccessMapperImpl{}
}

func (m *ChargeDataAccessMapperImpl) Entity(charge *domain.Charge) (*ChargeEntity, error) {
	var chargeID primitive.ObjectID

	if charge.ChargeID.String() != "" {
		objectID, err := primitive.ObjectIDFromHex(charge.ChargeID.String())
		if err != nil {
			return nil, fmt.Errorf("invalid ChargeID: %v", err)
		}

		chargeID = objectID
	}

	return &ChargeEntity{
		ChargeID:          chargeID,
		UserID:            charge.UserID.String(),
		ServiceID:         charge.ServiceID.String(),
		Timestamp:         charge.Timestamp.Int64(),
		AccumulatedAmount: charge.AccumulatedAmount.Int64(),
	}, nil
}

func (m *ChargeDataAccessMapperImpl) Domain(entity *ChargeEntity) *domain.Charge {
	return &domain.Charge{
		ChargeID:          valueobject.NewChargeID(entity.ChargeID.Hex()),
		UserID:            valueobject.NewUserID(entity.UserID),
		ServiceID:         valueobject.NewServiceID(entity.ServiceID),
		Timestamp:         valueobject.NewTimestamp(entity.Timestamp),
		AccumulatedAmount: valueobject.NewAmount(entity.AccumulatedAmount),
	}
}
