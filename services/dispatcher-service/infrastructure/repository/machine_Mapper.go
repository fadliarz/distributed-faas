package repository

import (
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/domain-core"
)

type MachineDataAccessMapperImpl struct{}

func NewMachineDataAccessMapper() MachineDataAccessMapper {
	return &MachineDataAccessMapperImpl{}
}

func (m *MachineDataAccessMapperImpl) Domain(entity *MachineEntity) *domain.Machine {
	if entity == nil {
		return nil
	}

	return &domain.Machine{
		MachineID: domain.NewMachineID(entity.MachineID.Hex()),
		Address:   domain.NewAddress(entity.Address),
		Status:    domain.NewStatus(entity.Status),
	}
}
