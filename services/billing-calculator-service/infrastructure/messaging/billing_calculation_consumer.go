package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/billing-calculator-service/domain/application-service"
	"github.com/rs/zerolog/log"
)

// Consumer

type BillingCalculationConsumerImpl struct {
	consumer *kafka.ConfluentKafkaConsumer[*application.BillingCalculationEvent]
}

func NewBillingCalculationConsumer(consumer *kafka.ConfluentKafkaConsumer[*application.BillingCalculationEvent]) application.BillingCalculationConsumer {
	return &BillingCalculationConsumerImpl{
		consumer: consumer,
	}
}

func (c *BillingCalculationConsumerImpl) PollAndProcessMessages() {
	c.consumer.PollAndProcessMessages()
}

// Deserializer

type BillingCalculationMessageDeserializer struct{}

func NewBillingCalculationMessageDeserializer() kafka.MessageDeserializer[*application.BillingCalculationEvent] {
	return &BillingCalculationMessageDeserializer{}
}

func (d *BillingCalculationMessageDeserializer) Deserialize(data []byte) (*application.BillingCalculationEvent, error) {
	var event application.BillingCalculationEvent
	if err := json.Unmarshal(data, &event); err != nil {
		log.Error().
			Err(err).
			Bytes("data", data).
			Msg("Failed to deserialize billing calculation event")

		return nil, fmt.Errorf("failed to deserialize billing calculation event: %w", err)
	}

	log.Debug().
		Str("user_id", event.UserID).
		Int64("last_billed", event.LastBilled).
		Msg("Successfully deserialized billing calculation event")

	return &event, nil
}

// Processor

type BillingCalculationMessageProcessor struct {
	handler application.BillingCalculationEventHandler
}

func NewBillingCalculationMessageProcessor(handler application.BillingCalculationEventHandler) kafka.MessageProcessor[*application.BillingCalculationEvent] {
	return &BillingCalculationMessageProcessor{
		handler: handler,
	}
}

func (p *BillingCalculationMessageProcessor) Process(ctx context.Context, event *application.BillingCalculationEvent) error {
	log.Info().
		Str("user_id", event.UserID).
		Int64("last_billed", event.LastBilled).
		Msg("Processing billing calculation event")

	err := p.handler.Handle(ctx, event)
	if err != nil {
		return err
	}

	return nil
}
