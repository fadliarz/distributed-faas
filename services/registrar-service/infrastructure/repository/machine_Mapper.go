package repository

import (
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MachineMapperImpl struct{}

func NewMachineMapper() MachineMapper {
	return &MachineMapperImpl{}
}

func (m *MachineMapperImpl) Entity(machine *domain.Machine) (*MachineEntity, error) {
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

func (m *MachineMapperImpl) Domain(entity *MachineEntity) *domain.Machine {
	if entity == nil {
		return nil
	}

	return &domain.Machine{
		MachineID: domain.NewMachineID(entity.MachineID.Hex()),
		Address:   domain.NewAddress(entity.Address),
		Status:    domain.NewStatus(entity.Status),
	}
}
