package repository

import (
	"context"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
)

type ChargeRepositoryImpl struct {
	mapper ChargeDataAccessMapper
	repo   *ChargeMongoRepository
}

func NewChargeRepository(mapper ChargeDataAccessMapper, repo *ChargeMongoRepository) application.ChargeRepository {
	return &ChargeRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *ChargeRepositoryImpl) FindChargesByUserIDAndTimeRange(ctx context.Context, userID valueobject.UserID, startTime, endTime int64) ([]domain.Charge, error) {
	entities, err := r.repo.FindChargesByUserIDAndTimeRange(ctx, userID.String(), startTime, endTime)
	if err != nil {
		return nil, err
	}

	charges := make([]domain.Charge, 0, len(entities))
	for _, entity := range entities {
		charge := r.mapper.Domain(entity)
		if charge != nil {
			charges = append(charges, *charge)
		}
	}

	return charges, nil
}
