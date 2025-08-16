package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/billing-cron-service/domain/domain-core"
)

// Constructs

type BillingCronApplicationService struct {
	repositoryManager *BillingCronApplicationServiceRepositoryManager
}

type BillingCronApplicationServiceRepositoryManager struct {
	Cron CronRepository
}

func NewBillingCronApplicationService(repositoryManager BillingCronApplicationServiceRepositoryManager) *BillingCronApplicationService {
	return &BillingCronApplicationService{
		repositoryManager: &repositoryManager,
	}
}

func NewBillingCronApplicationServiceRepositoryManager(cronRepository CronRepository) *BillingCronApplicationServiceRepositoryManager {
	return &BillingCronApplicationServiceRepositoryManager{
		Cron: cronRepository,
	}
}

// Methods

func (s *BillingCronApplicationService) UpdateLastBilled(ctx context.Context, timestampPair domain.TimestampPair) error {
	return s.repositoryManager.Cron.UpdateLastBilled(ctx, timestampPair)
}
