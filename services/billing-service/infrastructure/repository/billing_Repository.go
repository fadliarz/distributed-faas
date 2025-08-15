package repository

import (
	"context"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-service/domain/domain-core"
)

type BillingRepositoryImpl struct {
	mapper BillingDataAccessMapper
	repo   *BillingMongoRepository
}

func NewBillingRepository(mapper BillingDataAccessMapper, repo *BillingMongoRepository) application.BillingRepository {
	return &BillingRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *BillingRepositoryImpl) FindByUserID(ctx context.Context, userID valueobject.UserID) (*domain.Billing, error) {
	entity, err := r.repo.FindByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	return r.mapper.Domain(entity), nil
}
