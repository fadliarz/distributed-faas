package application

import (
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

type RegistrarDataMapperImpl struct{}

func NewRegistrarDataMapper() RegistrarDataMapper {
	return &RegistrarDataMapperImpl{}
}

func (m *RegistrarDataMapperImpl) CreateMachineCommandToMachine(command *CreateMachineCommand) *domain.Machine {
	return &domain.Machine{
		Address: domain.NewAddress(command.Address),
	}
}
