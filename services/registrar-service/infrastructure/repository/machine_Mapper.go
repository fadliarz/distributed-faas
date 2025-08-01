package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MachineDataAccessMapper struct{}

func NewMachineDataAccessMapper() MachineMapper {
	return &MachineDataAccessMapper{}
}

func (m *MachineDataAccessMapper) Entity(machine *domain.Machine) (*MachineEntity, error) {
	if machine == nil {
		return nil, nil
	}

	machineID, err := primitive.ObjectIDFromHex(machine.MachineID.String())
	if err != nil {
		return nil, fmt.Errorf("invalid machine ID: %w", err)
	}

	return &MachineEntity{
		MachineID: machineID,
		Address:   machine.Address.String(),
		Status:    machine.Status.String(),
	}, nil
}

func (m *MachineDataAccessMapper) Domain(entity *MachineEntity) *domain.Machine {
	if entity == nil {
		return nil
	}

	return &domain.Machine{
		MachineID: domain.NewMachineID(entity.MachineID.Hex()),
		Address:   domain.NewAddress(entity.Address),
		Status:    domain.NewStatus(entity.Status),
	}
}
