package handler

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/features/command"
	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

func (h *CommandHandler) CreateMachine(ctx context.Context, cmd *command.CreateMachineCommand) (domain.MachineID, error) {
	machineID, err := h.service.PersistMachine(ctx, cmd)
	if err != nil {
		return "", err
	}

	return machineID, nil
}
