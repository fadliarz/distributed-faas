package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/billing-cron-service/domain/domain-core"
)

type BillingCronHandler struct {
	service *BillingCronApplicationService
}

func NewBillingCronHandler(service *BillingCronApplicationService) *BillingCronHandler {
	return &BillingCronHandler{
		service: service,
	}
}

func (h *BillingCronHandler) UpdateLastBilled(ctx context.Context) error {
	timestampPair := domain.NewTimestampPair()

	err := h.service.UpdateLastBilled(ctx, timestampPair)
	if err != nil {
		return fmt.Errorf("failed to update LastBilled: %w", err)
	}

	return nil
}
