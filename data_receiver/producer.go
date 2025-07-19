package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(*types.OBUData) error
	Stop()
}

type KafkaDataProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaDataProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error == nil {
					log.Printf("Delivered to %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()
	return &KafkaDataProducer{
		producer: p,
		topic:    topic,
	}, nil
}

func (p *KafkaDataProducer) ProduceData(data *types.OBUData) error {
	b, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal OBUData: %w", err)
	}
	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	return nil
}

func (p *KafkaDataProducer) Stop() {
	p.producer.Flush(15_000)
	p.producer.Close()
}
