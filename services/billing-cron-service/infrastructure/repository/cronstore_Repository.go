package repository

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/billing-cron-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/billing-cron-service/domain/domain-core"
)

type CronRepositoryImpl struct {
	repo *CronMongoRepository
}

func NewCronRepository(repo *CronMongoRepository) application.CronRepository {
	return CronRepositoryImpl{
		repo: repo,
	}
}

func (c CronRepositoryImpl) UpdateLastBilled(ctx context.Context, timestampPair domain.TimestampPair) error {
	return c.repo.UpdateLastBilled(ctx, timestampPair.BeforeTimestamp(), timestampPair.AfterTimestamp())
}
