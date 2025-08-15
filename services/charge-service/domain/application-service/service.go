package application

import (
	"context"
)

type ChargeApplicationServiceImpl struct {
	dataMapper ChargeDataMapper
	aggregator ChargeAggregator
}

func NewChargeApplicationService(dataMapper ChargeDataMapper, aggregator ChargeAggregator) ChargeApplicationService {
	return &ChargeApplicationServiceImpl{
		dataMapper: dataMapper,
		aggregator: aggregator,
	}
}

func (s *ChargeApplicationServiceImpl) ProcessCharge(ctx context.Context, command *CreateChargeCommand) error {
	charge, err := s.dataMapper.CreateChargeCommandToCharge(command)
	if err != nil {
		return err
	}

	return s.aggregator.AddCharge(ctx, charge)
}
