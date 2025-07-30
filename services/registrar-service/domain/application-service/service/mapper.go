package service

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"github.com/rs/zerolog/log"
)

type Mapper interface {
	CreateMachineCommandToMachine(cmd *command.CreateMachineCommand) (*domain.Machine, error)
}

type MapperImpl struct{}

func NewMapper() Mapper {
	return &MapperImpl{}
}

func (m *MapperImpl) CreateMachineCommandToMachine(cmd *command.CreateMachineCommand) (*domain.Machine, error) {
	log.Debug().Msg("Mapping CreateMachineCommand to Machine")
	machineID, err := domain.NewMachineID(cmd.MachineID)
	if err != nil {
		return nil, fmt.Errorf("failed to create machine ID: %w", err)
	}

	return &domain.Machine{
		MachineID: machineID,
	}, nil
}
