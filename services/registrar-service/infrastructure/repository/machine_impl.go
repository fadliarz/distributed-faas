package repository

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/ports"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

type MachineRepositoryImpl struct {
	mapper MachineMapper
	repo   *MachineMongoRepository
}

func NewMachineRepository(mapper MachineMapper, repo *MachineMongoRepository) ports.MachineRepository {
	return &MachineRepositoryImpl{
		mapper: mapper,
		repo:   repo,
	}
}

func (r *MachineRepositoryImpl) Save(ctx context.Context, machine *domain.Machine) (domain.MachineID, error) {
	entity := r.mapper.Entity(machine)

	machineID, err := r.repo.Save(ctx, entity)
	if err != nil {
		return "", fmt.Errorf("failed to save machine: %w", err)
	}

	return domain.NewLooseMachineID(machineID), nil
}

func (r *MachineRepositoryImpl) UpdateStatus(ctx context.Context, machineID domain.MachineID, address domain.Address, status domain.Status) error {
	err := r.repo.UpdateStatus(ctx, machineID.String(), address.String(), status.String())
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	return nil
}
