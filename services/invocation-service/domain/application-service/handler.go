package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/invocation-service/domain/domain-core"
)

type CommandHandler struct {
	service *InvocationApplicationService
}

func NewCommandHandler(service *InvocationApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

func (h *CommandHandler) CreateInvocation(ctx context.Context, cmd *CreateInvocationCommand) (domain.InvocationID, error) {
	invocationID, err := h.service.PersistInvocation(ctx, cmd)
	if err != nil {
		return "", err
	}

	return invocationID, nil
}

func (h *CommandHandler) GetInvocation(ctx context.Context, query *GetInvocationQuery) (*domain.Invocation, error) {
	invocation, err := h.service.GetInvocation(ctx, query)
	if err != nil {
		return nil, err
	}

	return invocation, nil
}

