package application

import "context"

type RetryHandler struct {
	service *RetryApplicationService
}

func NewRetryHandler(service *RetryApplicationService) *RetryHandler {
	return &RetryHandler{
		service: service,
	}
}

func (h *RetryHandler) RetryInvocations(ctx context.Context) {

}
