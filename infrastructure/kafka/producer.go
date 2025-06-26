package kafka

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde"
	"github.com/confluentinc/confluent-kafka-go/v2/schemaregistry/serde/protobuf"
	"google.golang.org/protobuf/proto"
)

type ProducerDataMapper[M any, P proto.Message] interface {
	Proto(msg M) P
}

type ConfluentKafkaProducer[M any, P proto.Message] struct {
	producer           *confluentkafka.Producer
	protobufSerializer *protobuf.Serializer
	topicName          string
	dataMapper         ProducerDataMapper[M, P]
	producerContext    context.Context
	cancelProducer     context.CancelFunc
}

func NewConfluentKafkaProducer[M any, P proto.Message](cfg *kafka.ConfigMap, srClient schemaregistry.Client, topicName string,
	dataMapper ProducerDataMapper[M, P]) (*ConfluentKafkaProducer[M, P], error) {
	producer, err := kafka.NewProducer(cfg)
	if err != nil {
		return nil, err
	}

	protobufSerializer, err := protobuf.NewSerializer(srClient, serde.ValueSerde, protobuf.NewSerializerConfig())
	if err != nil {
		producer.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	kp := &ConfluentKafkaProducer[M, P]{
		producer:           producer,
		protobufSerializer: protobufSerializer,
		topicName:          topicName,
		dataMapper:         dataMapper,
		producerContext:    ctx,
		cancelProducer:     cancel,
	}

	go kp.handleDeliveryReports()

	return kp, nil
}

func (p *ConfluentKafkaProducer[M, P]) Produce(ctx context.Context, key string, msg M) error {
	protoMsg := p.dataMapper.Proto(msg)

	payload, err := p.protobufSerializer.Serialize(p.topicName, protoMsg)
	if err != nil {
		log.Printf("Failed to serialize message: %v", err)
		return err
	}

	err = p.producer.Produce(&confluentkafka.Message{
		TopicPartition: confluentkafka.TopicPartition{Topic: &p.topicName, Partition: int32(confluentkafka.PartitionAny)},
		Key:            []byte(key),
		Value:          payload,
	}, nil)
	if err != nil {
		log.Printf("Failed to produce message: %v", err)
		return err
	}

	return nil
}

func (p *ConfluentKafkaProducer[M, P]) handleDeliveryReports() {
	for {
		select {
		case e := <-p.producer.Events():
			switch ev := e.(type) {
			case *confluentkafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("Failed to deliver message: %v", ev.TopicPartition.Error)
					// ToDo: implement specific error handling, e.g., retry logic, dead-letter queue
				}
			case confluentkafka.Error:
				log.Printf("Producer error: %v", ev)
				// ToDo: handle fatal producer errors, potentially close and reinitialize
			}
		case <-p.producerContext.Done():
			return
		}
	}
}

func (p *ConfluentKafkaProducer[M, P]) Close() {
	p.cancelProducer()

	// wait for any outstanding messages to be delivered
	remainingMsgs := p.producer.Flush(15 * 1000)
	if remainingMsgs > 0 {
		log.Print("Warning: There are still ", remainingMsgs, " messages left in the producer queue after closing.")
	}

	// close the producer and serializer
	p.producer.Close()
	p.protobufSerializer.Close()
	log.Println("Confluent Kafka producer closed successfully.")
}
