package application

import (
	"context"
)

type CommandHandler struct {
	service ChargeApplicationService
}

func NewCommandHandler(service ChargeApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

func (h *CommandHandler) CreateCharge(ctx context.Context, command *CreateChargeCommand) error {
	return h.service.ProcessCharge(ctx, command)
}
