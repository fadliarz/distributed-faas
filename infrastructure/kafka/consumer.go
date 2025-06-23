package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)

type MessageProcessor[M any, P proto.Message] interface {
	ProcessMesssage(msg M) error
}

type DataMapper[M any, P proto.Message] interface {
	Message(protoMsg P) M
}

type ConfluentKafkaConsumer[M any, P proto.Message] struct {
	consumer             *confluentkafka.Consumer
	protobufDeserializer *protobuf.Deserializer
	topicName            string
	numProcessingWorkers int
	processingWg         sync.WaitGroup
	messageChannel       chan M
	consumerContext      context.Context
	cancelConsumer       context.CancelFunc
	stopProcessing       chan struct{}
	protoMessageFactory  func() P
	dataMapper           DataMapper[M, P]
	messageProcessor     MessageProcessor[M, P]
}

func NewConfluentKafkaConsumer[M any, P proto.Message](cfg *confluentkafka.ConfigMap, srClient schemaregistry.Client, topicName string, numProcessingWorkers int, protoMessageFactory func() P, dataMapper DataMapper[M, P], messageProcessor MessageProcessor[M, P]) (*ConfluentKafkaConsumer[M, P], error) {
	consumer, err := confluentkafka.NewConsumer(cfg)
	if err != nil {
		return nil, err
	}

	consumer.SubscribeTopics([]string{topicName}, nil)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("Failed subscribing to kafka topic (%s): %v", topicName, err)
	}

	protobufDeserializer, err := protobuf.NewDeserializer(srClient, serde.ValueSerde, protobuf.NewDeserializerConfig())
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("Failed registering schema register: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	service := &ConfluentKafkaConsumer[M, P]{
		consumer:             consumer,
		protobufDeserializer: protobufDeserializer,
		topicName:            topicName,
		numProcessingWorkers: numProcessingWorkers,
		consumerContext:      ctx,
		cancelConsumer:       cancel,
		stopProcessing:       make(chan struct{}),
		protoMessageFactory:  protoMessageFactory,
		dataMapper:           dataMapper,
		messageProcessor:     messageProcessor,
	}

	return service, nil
}

func (c *ConfluentKafkaConsumer[M, P]) Consume(ctx context.Context) {

	c.startProcessingWorkers()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigchan:
			return
		case <-c.consumerContext.Done():
			return
		case <-ctx.Done():
			return
		default:
			event := c.consumer.Poll(100)
			if event == nil {
				return
			}

			switch assertedEvent := event.(type) {
			case *confluentkafka.Message:
				protoMsg := c.protoMessageFactory()
				c.protobufDeserializer.DeserializeInto(*assertedEvent.TopicPartition.Topic, assertedEvent.Value, protoMsg)

				select {
				case c.messageChannel <- c.dataMapper.Message(protoMsg):
				case <-c.consumerContext.Done():
					return
				case <-time.After(5 * time.Second):
					// ToDo: handle timeout (e.g., log it or retry)
				}
			case confluentkafka.Error:
				if assertedEvent.IsFatal() {
					c.cancelConsumer()
					return
				} else {
					// ToDo: handle non-fatal errors (e.g., log them)
				}
			default:
				// ToDo: handle other event types if necessary
			}
		}
	}
}

func (c *ConfluentKafkaConsumer[M, P]) startProcessingWorkers() {
	for i := 0; i < c.numProcessingWorkers; i++ {
		c.processingWg.Add(1)
		go func(workerID int) {
			defer c.processingWg.Done()
			for {
				select {
				case msg, ok := <-c.messageChannel:
					if !ok {
						return
					}

					err := c.messageProcessor.ProcessMesssage(msg)
					if err != nil {
						// ToDo: handle processing error (e.g., dlq, log it, or retry)
					} else {
						// ToDo: handle successful processing (e.g., ack the message)
						/*
							ToDo:
							1. If using manual commits, commit the offset
							for this message or coordinate with a commit manager.
						*/
					}
				case <-c.consumerContext.Done():
					return
				case <-c.stopProcessing:
					return
				}
			}
		}(i)
	}
}
