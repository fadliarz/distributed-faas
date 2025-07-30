package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/dispatcher-service/domain/application-service"
)

// Consumer

type InvocationConsumerImpl struct {
	consumer *kafka.ConfluentKafkaConsumer[*application.InvocationCreatedEvent]
}

func NewInvocationConsumer(consumer *kafka.ConfluentKafkaConsumer[*application.InvocationCreatedEvent]) application.InvocationConsumer {
	return &InvocationConsumerImpl{
		consumer: consumer,
	}
}

func (c *InvocationConsumerImpl) PollAndProcessMessages() {
	c.consumer.PollAndProcessMessages()
}

// Deserializer

type InvocationMessageDeserializer struct{}

func NewInvocationMessageDeserializer() kafka.MessageDeserializer[*application.InvocationCreatedEvent] {
	return &InvocationMessageDeserializer{}
}

func (d *InvocationMessageDeserializer) Deserialize(data []byte) (*application.InvocationCreatedEvent, error) {
	var msg application.InvocationCreatedEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Processor

type InvocationMessageProcessor struct {
	handler *application.InvocationEventHandler
}

func NewInvocationMessageProcessor(handler *application.InvocationEventHandler) kafka.MessageProcessor[*application.InvocationCreatedEvent] {
	return &InvocationMessageProcessor{
		handler: handler,
	}
}

func (p *InvocationMessageProcessor) Process(ctx context.Context, event *application.InvocationCreatedEvent) error {
	err := p.handler.HandleInvocationCreatedEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("failed to process invocation created event: %w", err)
	}

	return nil
}
