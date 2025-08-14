package application

import (
	"context"
	"fmt"
	"time"

	"github.com/fadliarz/distributed-faas/common/valueobject"
	"github.com/fadliarz/distributed-faas/services/accumulator-service/domain/domain-core"
	"github.com/rs/zerolog/log"
)

type ChargeApplicationServiceImpl struct {
	chargeRepository ChargeRepository
	domainService    domain.AccumulatorDomainService
}

func NewChargeApplicationService(
	chargeRepository ChargeRepository,
) ChargeApplicationService {
	return &ChargeApplicationServiceImpl{
		chargeRepository: chargeRepository,
		domainService:    domain.NewAccumulatorDomainService(),
	}
}

func (s *ChargeApplicationServiceImpl) ProcessChargeEventBatch(ctx context.Context, events []*ChargeEvent) error {
	// Base case
	if len(events) == 0 {
		return nil
	}

	log.Info().Int("batch_size", len(events)).Msg("Processing accumulation batch")

	// Group events by UserID and ServiceID to aggregate amounts
	chargeMap := make(map[string]*domain.Charge)

	for _, event := range events {
		key := fmt.Sprintf("%s:%s", event.UserID.String(), event.ServiceID.String())
		timestamp := time.Date(time.Now().Year(), time.Now().Month(), 0, 0, 0, 0, 0, time.UTC).Unix()

		if entry, exists := chargeMap[key]; exists {
			entry.AccumulatedAmount = valueobject.NewAmount(entry.AccumulatedAmount.Int64() + event.Amount.Int64())
		} else {
			chargeMap[key] = &domain.Charge{
				UserID:            event.UserID,
				ServiceID:         event.ServiceID,
				Timestamp:         valueobject.NewTimestamp(timestamp),
				AccumulatedAmount: event.Amount,
			}
		}
	}

	charges := make([]*domain.Charge, 0, len(chargeMap))
	for _, charge := range chargeMap {
		charges = append(charges, charge)
	}

	err := s.chargeRepository.UpsertCharges(ctx, charges)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to upsert charges in batch")
	}

	return nil
}
