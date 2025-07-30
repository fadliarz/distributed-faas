package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/ports"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

var (
	ErrUserNotAuthorized = errors.New("you're not authorized to perform this action")
)

type RegistrarApplicationService struct {
	mapper            Mapper
	domainSvc         domain.RegistrarDomainService
	repositoryManager *RepositoryManager
}

type RepositoryManager struct {
	Machine ports.MachineRepository
}

func NewRegistrarApplicationService(mapper Mapper, domainSvc domain.RegistrarDomainService, machineRepo ports.MachineRepository) *RegistrarApplicationService {
	return &RegistrarApplicationService{
		mapper:    mapper,
		domainSvc: domainSvc,
		repositoryManager: &RepositoryManager{
			Machine: machineRepo,
		},
	}
}

func (s *RegistrarApplicationService) PersistMachine(ctx context.Context, cmd *command.CreateMachineCommand) (domain.MachineID, error) {
	// Validate the command
	machine, err := s.mapper.CreateMachineCommandToMachine(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to map command to invocation: %w", err)
	}

	// Validate and initiate the machine
	err = s.domainSvc.ValidateAndInitiateMachine(machine)
	if err != nil {
		return "", fmt.Errorf("failed to validate and initiate machine: %w", err)
	}

	// Save the machine
	machineID, err := s.repositoryManager.Machine.Save(ctx, machine)
	if err != nil {
		return "", fmt.Errorf("failed to save machine: %w", err)
	}

	return machineID, nil
}

func (s *RegistrarApplicationService) UpdateMachineStatus(ctx context.Context, cmd *command.UpdateMachineStatusCommand) error {
	// Validate the command
	machineID, err := domain.NewMachineID(cmd.MachineID)
	if err != nil {
		return fmt.Errorf("invalid machine ID: %w", err)
	}

	address, err := domain.NewAddress(cmd.Address)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	// Update the machine status
	err = s.repositoryManager.Machine.UpdateStatus(ctx, machineID, address, domain.NewStatusFromInt(int(domain.Available)))
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	return nil
}
