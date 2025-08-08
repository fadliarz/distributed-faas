package application

import (
	"context"
	"fmt"

	"github.com/fadliarz/distributed-faas/services/retry-service/domain/domain-core"
)

type RetryHandler struct {
	service *RetryApplicationService
}

func NewRetryHandler(service *RetryApplicationService) *RetryHandler {
	return &RetryHandler{
		service: service,
	}
}

func (h *RetryHandler) RetryInvocations(ctx context.Context, thresholdInSec int64) error {
	err := h.service.RetryInvocations(ctx, domain.NewThreshold(thresholdInSec))
	if err != nil {
		return fmt.Errorf("failed to retry invocations: %w", err)
	}

	return nil
}
