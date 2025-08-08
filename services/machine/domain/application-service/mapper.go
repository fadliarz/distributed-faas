package application

import (
	"github.com/fadliarz/distributed-faas/services/machine/domain/domain-core"
)

type MachineDataMapperImpl struct{}

func NewMachineDataMapper() MachineDataMapper {
	return &MachineDataMapperImpl{}
}

func (m *MachineDataMapperImpl) ProcessInvocationCommandToCheckpoint(cmd *ProcessInvocationCommand) *domain.Checkpoint {
	return &domain.Checkpoint{
		CheckpointID:  domain.NewCheckpointID(cmd.InvocationID),
		FunctionID:    domain.NewFunctionID(cmd.FunctionID),
		UserID:        domain.NewUserID(cmd.UserID),
		SourceCodeURL: domain.NewSourceCodeURL(cmd.SourceCodeURL),
		Status:        domain.NewStatus(cmd.Status),
		Timestamp:     domain.NewTimestamp(cmd.Timestamp),
	}
}
