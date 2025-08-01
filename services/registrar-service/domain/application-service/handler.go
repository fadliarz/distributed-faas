package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/registrar-service/domain/domain-core"
)

type CommandHandler struct {
	service *RegistrarApplicationService
}

func NewCommandHandler(service *RegistrarApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

func (h *CommandHandler) CreateMachine(ctx context.Context, command *CreateMachineCommand) (*domain.Machine, error) {
	machine, err := h.service.PersistMachine(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to persist machine: %w", err)
	}

	return machine, nil
}

func (h *CommandHandler) UpdateMachineStatus(ctx context.Context, command *UpdateMachineStatusCommand) error {
	err := h.service.UpdateMachineStatus(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to update machine status: %w", err)
	}

	return nil
}
