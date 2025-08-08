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

func (h *CommandHandler) CreateInvocation(ctx context.Context, cmd *CreateInvocationCommand) (*domain.Invocation, error) {
	return h.service.PersistInvocation(ctx, cmd)
}

func (h *CommandHandler) GetInvocation(ctx context.Context, query *GetInvocationQuery) (*domain.Invocation, error) {
	return h.service.GetInvocation(ctx, query)
}
