package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/domain-core"
)

type MachineRepositoryImpl struct {
	mapper MachineDataAccessMapper
	repo   *MachineMongoRepository
}

func NewMachineRepository(mapper MachineDataAccessMapper, repo *MachineMongoRepository) application.MachineRepository {
	return &MachineRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *MachineRepositoryImpl) FindManyByStatus(ctx context.Context, status domain.Status) ([]domain.Machine, error) {
	machines, err := r.repo.FindManyByStatus(ctx, status.String())
	if err != nil {
		return nil, fmt.Errorf("failed to find machines by status: %w", err)
	}

	if machines == nil {
		return []domain.Machine{}, nil
	}

	var result []domain.Machine
	for _, machine := range machines {
		result = append(result, *r.mapper.Domain(&machine))
	}

	return result, nil
}
