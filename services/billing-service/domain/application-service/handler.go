package application

import (
	"context"

	"github.com/fadliarz/distributed-faas/services/billing-service/domain/domain-core"
)

// Constructor

type CommandHandler struct {
	service BillingApplicationService
}

func NewCommandHandler(service BillingApplicationService) *CommandHandler {
	return &CommandHandler{
		service: service,
	}
}

// Methods

func (h *CommandHandler) GetBilling(ctx context.Context, query *GetBillingQuery) (*domain.Billing, error) {
	return h.service.GetBilling(ctx, query)
}
