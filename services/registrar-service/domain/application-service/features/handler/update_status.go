package handler

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/application-service/features/command"
)

func (h *CommandHandler) UpdateMachineStatus(ctx context.Context, cmd *command.UpdateMachineStatusCommand) error {
	err := h.service.UpdateMachineStatus(ctx, cmd)
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	return nil
}
