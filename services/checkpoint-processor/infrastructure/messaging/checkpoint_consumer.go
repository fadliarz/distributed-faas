package messaging

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/fadliarz/distributed-faas/infrastructure/kafka"
	"github.com/fadliarz/distributed-faas/services/checkpoint-processor/domain/application-service"
	"github.com/rs/zerolog/log"
)

// Consumer

type CheckpointConsumerImpl struct {
	consumer *kafka.ConfluentKafkaConsumer[*application.CheckpointEvent]
}

func NewCheckpointConsumer(consumer *kafka.ConfluentKafkaConsumer[*application.CheckpointEvent]) application.CheckpointConsumer {
	return &CheckpointConsumerImpl{
		consumer: consumer,
	}
}
func (c *CheckpointConsumerImpl) PollAndProcessMessages() {
	c.consumer.PollAndProcessMessages()
}

// Deserializer

type CheckpointMessageDeserializer struct{}

func NewCheckpointMessageDeserializer() kafka.MessageDeserializer[*application.CheckpointEvent] {
	return &CheckpointMessageDeserializer{}
}

func (d *CheckpointMessageDeserializer) Deserialize(data []byte) (*application.CheckpointEvent, error) {
	var msg application.CheckpointEvent
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// Processor

type CheckpointMessageProcessor struct {
	handler *application.CheckpointEventHandler
}

func NewCheckpointMessageProcessor(handler *application.CheckpointEventHandler) kafka.MessageProcessor[*application.CheckpointEvent] {
	return &CheckpointMessageProcessor{
		handler: handler,
	}
}

func (p *CheckpointMessageProcessor) Process(ctx context.Context, msg *application.CheckpointEvent) error {
	log.Info().Msgf("Processing checkpoint event: id (%s), outputURL (%s)", msg.CheckpointID, msg.OutputURL)

	exception := p.handler.HandleCheckpointEvent(ctx, msg)
	if exception != nil {
		return fmt.Errorf("failed to handle checkpoint event: %w", exception)
	}

	return nil
}
