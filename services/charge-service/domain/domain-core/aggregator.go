package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type ChargeAggregatorImpl struct {
	domainService    ChargeDomainService
	producer         ChargeProducer
	aggregates       map[string]*ChargeAggregate
	aggregatesMutex  sync.RWMutex
	aggregationTimer *time.Timer
	duration         time.Duration
	isRunning        bool
	runningMutex     sync.RWMutex
	stopChan         chan struct{}
}

func NewChargeAggregator(domainService ChargeDomainService, producer ChargeProducer, duration time.Duration) ChargeAggregator {
	return &ChargeAggregatorImpl{
		domainService: domainService,
		producer:      producer,
		aggregates:    make(map[string]*ChargeAggregate),
		duration:      duration,
		stopChan:      make(chan struct{}),
	}
}

func (a *ChargeAggregatorImpl) AddCharge(ctx context.Context, charge *Charge) error {
	if !a.isStarted() {
		return fmt.Errorf("charge aggregator is not started")
	}

	if err := a.domainService.ValidateCharge(charge); err != nil {
		return err
	}

	a.aggregatesMutex.Lock()
	defer a.aggregatesMutex.Unlock()

	key := charge.UserID.String() + ":" + charge.ServiceID.String()
	aggregate, exists := a.aggregates[key]

	if !exists {
		aggregate = NewChargeAggregate(charge.UserID, charge.ServiceID)
		a.aggregates[key] = aggregate
	}

	aggregate.AggregatedAmount = aggregate.AggregatedAmount.Add(charge.Amount)

	return nil
}

func (a *ChargeAggregatorImpl) Start(ctx context.Context) {
	a.runningMutex.Lock()
	if a.isRunning {
		a.runningMutex.Unlock()
		return
	}
	a.isRunning = true
	a.runningMutex.Unlock()

	a.resetTimer()

	for {
		select {
		case <-a.stopChan:
			a.runningMutex.Lock()
			a.isRunning = false
			a.runningMutex.Unlock()
			return
		case <-ctx.Done():
			a.runningMutex.Lock()
			a.isRunning = false
			a.runningMutex.Unlock()
			return
		case <-a.aggregationTimer.C:
			a.flushAggregates(ctx)
			a.resetTimer()
		}
	}
}

func (a *ChargeAggregatorImpl) Stop() {
	if a.isStarted() {
		close(a.stopChan)
	}
}

func (a *ChargeAggregatorImpl) isStarted() bool {
	a.runningMutex.RLock()

	defer a.runningMutex.RUnlock()

	return a.isRunning
}

func (a *ChargeAggregatorImpl) resetTimer() {
	if a.aggregationTimer != nil {
		a.aggregationTimer.Stop()
	}

	a.aggregationTimer = time.NewTimer(a.duration)
}

func (a *ChargeAggregatorImpl) flushAggregates(ctx context.Context) {
	a.aggregatesMutex.Lock()

	if len(a.aggregates) == 0 {
		a.aggregatesMutex.Unlock()

		log.Debug().Msg("No aggregates to flush")

		return
	}

	// Create a slice of aggregates to send
	aggregatesToSend := make([]*ChargeAggregate, 0, len(a.aggregates))
	for _, aggregate := range a.aggregates {
		aggregatesToSend = append(aggregatesToSend, aggregate)
	}

	// Clear the aggregates map
	a.aggregates = make(map[string]*ChargeAggregate)
	a.aggregatesMutex.Unlock()

	log.Info().Int("count", len(aggregatesToSend)).Msg("Flushing charge aggregates")

	// Send aggregates to Kafka
	if err := a.producer.SendAggregatedCharges(ctx, aggregatesToSend); err != nil {
		log.Error().Err(err).Msg("Failed to send aggregated charges")

		// TODO: Implement retry logic or dead letter queue
	}
}
