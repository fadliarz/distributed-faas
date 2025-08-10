package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/user-processor/domain/application-service"
	"github.com/rs/zerolog/log"
)

// Consumer

type UserConsumerImpl struct {
	consumer *kafka.ConfluentKafkaConsumer[*application.UserEvent]
}

func NewUserConsumer(consumer *kafka.ConfluentKafkaConsumer[*application.UserEvent]) application.UserConsumer {
	return &UserConsumerImpl{
		consumer: consumer,
	}
}

func (c *UserConsumerImpl) PollAndProcessMessages() {
	c.consumer.PollAndProcessMessages()
}

// Deserializer

type UserMessageDeserializer struct{}

func NewUserMessageDeserializer() kafka.MessageDeserializer[*application.UserEvent] {
	return &UserMessageDeserializer{}
}

func (d *UserMessageDeserializer) Deserialize(data []byte) (*application.UserEvent, error) {
	var msg application.UserEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Processor

type UserMessageProcessor struct {
	handler *application.UserEventHandler
}

func NewUserMessageProcessor(handler *application.UserEventHandler) kafka.MessageProcessor[*application.UserEvent] {
	return &UserMessageProcessor{
		handler: handler,
	}
}

func (p *UserMessageProcessor) Process(ctx context.Context, msg *application.UserEvent) error {
	log.Info().Msgf("Processing user event: userID (%s)", msg.UserID)

	exception := p.handler.HandleUserEvent(ctx, msg)
	if exception != nil {
		return fmt.Errorf("failed to handle user event: %w", exception)
	}

	return nil
}
