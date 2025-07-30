package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/function-service/domain/domain-core"
)

// Constructor

type CommandHandler struct {
	service FunctionApplicationService
}

func NewCommandHandler(service FunctionApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

// Methods

func (h *CommandHandler) CreateFunction(ctx context.Context, command *CreateFunctionCommand) (domain.FunctionID, error) {
	functionID, err := h.service.PersistFunction(ctx, command)
	if err != nil {
		return "", fmt.Errorf("failed to create function: %w", err)
	}

	return functionID, nil
}

func (h *CommandHandler) GetPresignedURL(ctx context.Context) (string, error) {
	

	return "", nil
}
