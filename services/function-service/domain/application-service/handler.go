package application

import (
	"context"

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
	return h.service.PersistFunction(ctx, command)
}

func (h *CommandHandler) GetFunctionUploadPresignedURL(ctx context.Context, query *GetFunctionUploadPresignedURLQuery) (string, error) {
	return h.service.GetFunctionUploadPresignedURL(ctx, query)
}

func (h *CommandHandler) UpdateFunctionSourceCodeURL(ctx context.Context, command *UpdateFunctionSourceCodeURLCommand) error {
	return h.service.UpdateFunctionSourceCodeURL(ctx, command)
}
