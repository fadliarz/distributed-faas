package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

// Ports

type MachineRepository interface {
	Save(ctx context.Context, machine *domain.Machine) (domain.MachineID, error)
	UpdateStatus(ctx context.Context, machineID domain.MachineID, address domain.Address, status domain.Status) error
}

// Interfaces

type RegistrarDataMapper interface {
	CreateMachineCommandToMachine(command *CreateMachineCommand) *domain.Machine
}
