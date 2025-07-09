package main

import (
	"encoding/json"
	"fmt"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type DataProducer interface {
	ProduceData(*types.OBUData) error
}

type KafkaDataProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaDataProducer(topic string) (DataProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, err
	}
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error == nil {
					fmt.Printf("Delivered to %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
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
		return err
	}
	return p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &p.topic,
			Partition: kafka.PartitionAny,
		},
		Value: b,
	}, nil)
}
