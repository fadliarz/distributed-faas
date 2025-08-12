package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
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

func (r *BillingRepositoryImpl) Save(ctx context.Context, billing *domain.Billing) (valueobject.BillingID, error) {
	entity, err := r.mapper.Entity(billing)
	if err != nil {
		return "", fmt.Errorf("failed to map billing entity: %w", err)
	}

	billingID, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("failed to save billing: %w", err)
	}

	return valueobject.NewBillingID(billingID), nil
}
