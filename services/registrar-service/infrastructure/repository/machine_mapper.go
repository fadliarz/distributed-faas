package repository

import (
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MachineMapperImpl struct{}

func NewMachineMapper() MachineMapper {
	return &MachineMapperImpl{}
}

func (m *MachineMapperImpl) Entity(machine *domain.Machine) *MachineEntity {
	if machine == nil {
		return nil
	}

	machineID := primitive.NilObjectID
	if machine.MachineID.String() != "" {
		if objectID, err := primitive.ObjectIDFromHex(machine.MachineID.String()); err == nil {
			machineID = objectID
		}
	}

	return &MachineEntity{
		MachineID: machineID,
		Address:   machine.Address.String(),
		Status:    machine.Status.String(),
	}
}

func (m *MachineMapperImpl) Domain(entity *MachineEntity) *domain.Machine {
	if entity == nil {
		return nil
	}

	return &domain.Machine{
		MachineID: domain.NewLooseMachineID(entity.MachineID.Hex()),
		Address:   domain.NewLooseAddress(entity.Address),
		Status:    domain.NewStatus(entity.Status),
	}
}
