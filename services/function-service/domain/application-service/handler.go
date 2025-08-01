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

func (h *CommandHandler) CreateFunction(ctx context.Context, command *CreateFunctionCommand) (*domain.Function, error) {
	functionID, err := h.service.PersistFunction(ctx, command)
	if err != nil {
		return nil, fmt.Errorf("failed to create function: %w", err)
	}

	return functionID, nil

	return functionID, nil
}

func (h *CommandHandler) GetFunctionUploadPresignedURL(ctx context.Context, query *GetFunctionUploadPresignedURLQuery) (string, error) {
	url, err := h.service.GetFunctionUploadPresignedURL(ctx, query)
	if err != nil {
		return "", fmt.Errorf("failed to get presigned URL for function upload: %w", err)
	}

	return url, nil
}

func (h *CommandHandler) UpdateFunctionSourceCodeURL(ctx context.Context, command *UpdateFunctionSourceCodeURLCommand) error {
	err := h.service.UpdateFunctionSourceCodeURL(ctx, command)
	if err != nil {
		return fmt.Errorf("failed to update function source code URL: %w", err)
	}

	return nil
}
