package application

import (
	"context"
	"fmt"
	"time"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/domain-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillingCalculatorApplicationServiceImpl struct {
	mapper            BillingCalculatorDataMapper
	domainService     domain.BillingCalculatorDomainService
	repositoryManager *BillingCalculatorApplicationServiceRepositoryManager
}

type BillingCalculatorApplicationServiceRepositoryManager struct {
	Charge  ChargeRepository
	Billing BillingRepository
}

func NewBillingCalculatorApplicationService(mapper BillingCalculatorDataMapper, domainService domain.BillingCalculatorDomainService, repositoryManager *BillingCalculatorApplicationServiceRepositoryManager) BillingCalculatorApplicationService {
	return &BillingCalculatorApplicationServiceImpl{
		mapper:            mapper,
		domainService:     domainService,
		repositoryManager: repositoryManager,
	}
}

func NewBillingCalculatorApplicationServiceRepositoryManager(chargeRepository ChargeRepository, billingRepository BillingRepository) *BillingCalculatorApplicationServiceRepositoryManager {
	return &BillingCalculatorApplicationServiceRepositoryManager{
		Charge:  chargeRepository,
		Billing: billingRepository,
	}
}

func (s *BillingCalculatorApplicationServiceImpl) ProcessBillingCalculation(ctx context.Context, event *BillingCalculationEvent) error {
	billing := s.mapper.BillingCalculationEventToBilling(event)

	charges, err := s.repositoryManager.Charge.FindChargesByUserIDAndTimeRange(ctx, billing.UserID, valueobject.NewTimestamp(event.LastBilled))
	if err != nil {
		return fmt.Errorf("failed to retrieve charges: %w", err)
	}

	s.domainService.CalculateBilling(billing, charges)

	billingID := valueobject.NewBillingID(primitive.NewObjectIDFromTimestamp(time.Unix(event.LastBilled, 0)).Hex())

	err = s.domainService.ValidateAndInitiateBilling(billing, billingID)
	if err != nil {
		return fmt.Errorf("failed to validate and initiate billing: %w", err)
	}

	billingID, err = s.repositoryManager.Billing.Save(ctx, billing)
	if err != nil {
		return fmt.Errorf("failed to save billing record: %w", err)
	}

	return nil
}
