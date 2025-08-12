package application

import "context"

type BillingCalculationEventHandlerImpl struct {
	applicationService BillingCalculatorApplicationService
}

func NewBillingCalculationEventHandler(applicationService BillingCalculatorApplicationService) BillingCalculationEventHandler {
	return &BillingCalculationEventHandlerImpl{
		applicationService: applicationService,
	}
}

func (h *BillingCalculationEventHandlerImpl) Handle(ctx context.Context, event *BillingCalculationEvent) error {
	return h.applicationService.ProcessBillingCalculation(ctx, event)
}
