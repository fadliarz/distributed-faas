package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/domain-core"
)

type ChargeRepositoryImpl struct {
	mapper ChargeDataAccessMapper
	repo   *ChargeMongoRepository
}

func NewChargeRepository(
	mapper ChargeDataAccessMapper,
	repo *ChargeMongoRepository,
) application.ChargeRepository {
	return &ChargeRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *ChargeRepositoryImpl) UpsertCharges(ctx context.Context, charges []*domain.Charge) error {
	if len(charges) == 0 {
		return nil
	}

	var entities []*ChargeEntity
	for _, charge := range charges {
		entity, err := r.mapper.Entity(charge)
		if err != nil {
			return fmt.Errorf("failed to map charge to entity: %v", err)
		}

		entities = append(entities, entity)
	}

	return r.repo.UpsertCharges(ctx, entities)
}
