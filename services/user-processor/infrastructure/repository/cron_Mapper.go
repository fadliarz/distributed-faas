package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	domain "github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CronDataAccessMapperImpl struct{}

func NewCronDataAccessMapper() CronDataAccessMapper {
	return &CronDataAccessMapperImpl{}
}

func (m *CronDataAccessMapperImpl) Entity(cron *domain.Cron) (*CronEntity, error) {
	if cron == nil {
		return nil, nil
	}

	primitiveUserID, err := primitive.ObjectIDFromHex(cron.UserID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	return &CronEntity{
		UserID:      primitiveUserID.Hex(),
		LastBilling: cron.LastBilling.Int64(),
	}, nil
}

func (m *CronDataAccessMapperImpl) Domain(entity *CronEntity) *domain.Cron {
	if entity == nil {
		return nil
	}

	return &domain.Cron{
		UserID:      valueobject.NewUserID(entity.UserID),
		LastBilling: valueobject.NewLastBilled(entity.LastBilling),
	}
}
