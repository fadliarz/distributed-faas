package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

type MachineRepositoryImpl struct {
	mapper MachineMapper
	repo   *MachineMongoRepository
}

func NewMachineRepository(mapper MachineMapper, repo *MachineMongoRepository) application.MachineRepository {
	return &MachineRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *MachineRepositoryImpl) Save(ctx context.Context, machine *domain.Machine) (domain.MachineID, error) {
	entity, err := r.mapper.Entity(machine)
	if err != nil {
		return "", fmt.Errorf("failed to map machine to entity: %w", err)
	}

	machineID, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("failed to save machine: %w", err)
	}

	return domain.NewMachineID(machineID), nil
}

func (r *MachineRepositoryImpl) UpdateStatus(ctx context.Context, machineID domain.MachineID, address domain.Address, status domain.Status) error {
	err := r.repo.UpdateStatus(ctx, machineID.String(), address.String(), status.String())
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	return nil
}
