package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/user-processor/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/domain-core"
)

type CronRepositoryImpl struct {
	mapper CronDataAccessMapper
	repo   *CronMongoRepository
}

func NewCronRepository(mapper CronDataAccessMapper, repo *CronMongoRepository) application.CronRepository {
	return &CronRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *CronRepositoryImpl) Save(ctx context.Context, cron *domain.Cron) (domain.CronID, error) {
	entity, err := r.mapper.Entity(cron)
	if err != nil {
		return "", fmt.Errorf("failed to map cron entity: %w", err)
	}

	cronID, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", err
	}

	return domain.NewCronID(cronID), nil
}
