package application

import (
	"context"
)

type ChargeEventHandler struct {
	service ChargeApplicationService
}

func NewChargeEventHandler(service ChargeApplicationService) *ChargeEventHandler {
	return &ChargeEventHandler{
		service: service,
	}
}

func (h *ChargeEventHandler) HandleChargeEventBatch(ctx context.Context, events []*ChargeEvent) error {
	return h.service.ProcessChargeEventBatch(ctx, events)
}
