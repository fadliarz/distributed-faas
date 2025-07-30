package kafka

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"

	"github.com/rs/zerolog/log"
)

type MessageWithOffset[M any] struct {
	Message      M
	KafkaMessage *kafka.Message
}

type ConfluentKafkaConsumer[M any] struct {
	ctx    context.Context
	config *ConsumerConfig

	consumer     *confluentkafka.Consumer
	deserializer MessageDeserializer[M]

	processor      MessageProcessor[M]
	processingWg   sync.WaitGroup
	messageChannel chan MessageWithOffset[M]

	commitTicker *time.Ticker
}

func NewConfluentKafkaConsumer[M any](ctx context.Context, config *ConsumerConfig, deserializer MessageDeserializer[M], processor MessageProcessor[M]) (*ConfluentKafkaConsumer[M], error) {
	consumer, err := confluentkafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  config.Basic.BootstrapServers,
		"group.id":           config.Basic.GroupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, err
	}

	if err := consumer.SubscribeTopics([]string{config.Basic.Topic}, nil); err != nil {
		consumer.Close()
		return nil, fmt.Errorf("Failed subscribing to kafka topic (%s): %v", config.Basic.Topic, err)
	}

	service := &ConfluentKafkaConsumer[M]{
		ctx:            ctx,
		config:         config,
		consumer:       consumer,
		deserializer:   deserializer,
		processor:      processor,
		processingWg:   sync.WaitGroup{},
		messageChannel: make(chan MessageWithOffset[M], config.Processing.NumWorkers*2), // Buffered channel to hold messages for processing
		commitTicker:   time.NewTicker(5 * time.Second),
	}

	go service.commitManager()

	return service, nil
}

func (c *ConfluentKafkaConsumer[M]) PollAndProcessMessages() {
	log.Info().Msgf("Starting Kafka consumer for topic %s with group ID %s", c.config.Basic.Topic, c.config.Basic.GroupID)

	c.startProcessingWorkers()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-sigchan:
			return
		case <-c.ctx.Done():
			return
		default:
			msg, err := c.consumer.ReadMessage(c.config.Basic.PollTimeout)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.IsTimeout() {
					log.Debug().Msgf("No messages received within timeout period (%s)", c.config.Basic.PollTimeout)
				} else {
					log.Warn().Msgf("Error reading message: %v", err)
				}
				continue
			}

			log.Debug().Msgf("Received message: %s", msg.Value)

			if deserializedMsg, err := c.deserializer.Deserialize(msg.Value); err == nil {
				c.messageChannel <- MessageWithOffset[M]{Message: deserializedMsg, KafkaMessage: msg}
			} else {
				log.Warn().Msgf("Failed to deserialize message: %v", err)
			}
		}
	}
}

func (c *ConfluentKafkaConsumer[M]) startProcessingWorkers() {
	for i := 0; i < c.config.Processing.NumWorkers; i++ {
		c.processingWg.Add(1)
		go func(workerID int) {
			defer c.processingWg.Done()
			for {
				select {
				case <-c.ctx.Done():
					return
				case messageWithOffset, ok := <-c.messageChannel:
					if !ok {
						return
					}

					err := c.processor.Process(c.ctx, messageWithOffset.Message)
					if err != nil {
						log.Error().Msgf("Error processing message in worker %d: %v", workerID, err)
						// ToDo: handle processing error (e.g., dlq, log it, or retry)
					} else {
						log.Info().Msgf("Worker %d successfully processed message: %v", workerID, messageWithOffset.Message)

						for i := 0; i < 5; i++ {
							_, err = c.consumer.StoreMessage(messageWithOffset.KafkaMessage)

							if err == nil {
								return
							}

							log.Warn().Msgf("Failed to store message offset after processing: %v, retrying (%d/5)", err, i+1)

							time.Sleep(2 * time.Second)
						}

						log.Fatal().Err(fmt.Errorf("Failed to store message offset after processing: %v", err))
					}
				}
			}
		}(i)
	}
}

func (c *ConfluentKafkaConsumer[M]) commitManager() {
	for {
		select {
		case <-c.ctx.Done():
			c.commitOffsets()
			return
		case <-c.commitTicker.C:
			c.commitOffsets()
		}
	}
}

func (c *ConfluentKafkaConsumer[M]) commitOffsets() {
	var err error

	for i := 0; i < 10; i++ {
		_, err = c.consumer.Commit()

		if err == nil {
			log.Debug().Msg("Offsets committed successfully")

			return
		}

		if err.(confluentkafka.Error).Code() == confluentkafka.ErrNoOffset {
			log.Debug().Msg("No offsets to commit, skipping")

			return
		}

		log.Warn().Msgf("Failed to commit offsets: %v, retrying (%d/10)", err, i+1)

		time.Sleep(2 * time.Second)
	}

	log.Fatal().Err(fmt.Errorf("Failed to commit offsets after multiple attempts: %v", err))
}
