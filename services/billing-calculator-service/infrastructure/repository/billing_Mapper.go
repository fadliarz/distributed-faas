package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillingDataAccessMapperImpl struct{}

func NewBillingDataAccessMapper() BillingDataAccessMapper {
	return &BillingDataAccessMapperImpl{}
}

func (m *BillingDataAccessMapperImpl) Entity(billing *domain.Billing) (*BillingEntity, error) {
	if billing == nil {
		return nil, nil
	}

	billingID, err := primitive.ObjectIDFromHex(billing.BillingID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid billing ID: %w", err)
	}

	return &BillingEntity{
		BillingID: billingID,
		UserID:    billing.UserID.String(),
		Amount:    billing.Amount.Int64(),
	}, nil
}

func (m *BillingDataAccessMapperImpl) Domain(entity *BillingEntity) *domain.Billing {
	if entity == nil {
		return nil
	}

	return &domain.Billing{
		BillingID: valueobject.NewBillingID(entity.BillingID.Hex()),
		UserID:    valueobject.NewUserID(entity.UserID),
		Amount:    valueobject.NewAmount(entity.Amount),
	}
}
