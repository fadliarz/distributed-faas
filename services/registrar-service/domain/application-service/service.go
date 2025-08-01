package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegistrarApplicationService struct {
	mapper            RegistrarDataMapper
	domainService     domain.RegistrarDomainService
	repositoryManager *RegistrarApplicationServiceRepositoryManager
}

type RegistrarApplicationServiceRepositoryManager struct {
	Machine MachineRepository
}

func NewRegistrarApplicationService(mapper RegistrarDataMapper, domainService domain.RegistrarDomainService, repositoryManager *RegistrarApplicationServiceRepositoryManager) *RegistrarApplicationService {
	return &RegistrarApplicationService{
		mapper:            mapper,
		domainService:     domainService,
		repositoryManager: repositoryManager,
	}
}

func NewRegistrarApplicationServiceRepositoryManager(machineRepository MachineRepository) *RegistrarApplicationServiceRepositoryManager {
	return &RegistrarApplicationServiceRepositoryManager{
		Machine: machineRepository,
	}
}

func (s *RegistrarApplicationService) PersistMachine(ctx context.Context, command *CreateMachineCommand) (*domain.Machine, error) {
	machine := s.mapper.CreateMachineCommandToMachine(command)

	err := s.domainService.ValidateAndInitiateMachine(machine, domain.NewMachineID(primitive.NewObjectID().Hex()))
	if err != nil {
		return nil, fmt.Errorf("failed to validate and initiate machine: %w", err)
	}

	_, err = s.repositoryManager.Machine.Save(ctx, machine)
	if err != nil {
		return nil, fmt.Errorf("failed to save machine: %w", err)
	}

	return machine, nil
}

func (s *RegistrarApplicationService) UpdateMachineStatus(ctx context.Context, command *UpdateMachineStatusCommand) error {
	err := s.repositoryManager.Machine.UpdateStatus(ctx, domain.NewMachineID(command.MachineID), domain.NewAddress(command.Address), domain.NewStatusFromInt(int(domain.Available)))
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	return nil
}
