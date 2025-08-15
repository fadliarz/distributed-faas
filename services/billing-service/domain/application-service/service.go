package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-service/domain/domain-core"
)

// Constructors

type BillingApplicationServiceImpl struct {
	repositoryManager *BillingApplicationServiceRepositoryManager
}

type BillingApplicationServiceRepositoryManager struct {
	Billing BillingRepository
}

func NewBillingApplicationService(repositoryManager *BillingApplicationServiceRepositoryManager) *BillingApplicationServiceImpl {
	return &BillingApplicationServiceImpl{
		repositoryManager: repositoryManager,
	}
}

func NewBillingApplicationServiceRepositoryManager(billing BillingRepository) *BillingApplicationServiceRepositoryManager {
	return &BillingApplicationServiceRepositoryManager{
		Billing: billing,
	}
}

// Methods

func (s *BillingApplicationServiceImpl) GetBilling(ctx context.Context, query *GetBillingQuery) (*domain.Billing, error) {
	billing, err := s.repositoryManager.Billing.FindByUserID(ctx, valueobject.NewUserID(query.UserID))
	if err != nil {
		return nil, fmt.Errorf("failed to find billing by user ID: %w", err)
	}

	if billing == nil {
		return nil, domain.NewErrBillingNotFound(nil)
	}

	return billing, nil
}
