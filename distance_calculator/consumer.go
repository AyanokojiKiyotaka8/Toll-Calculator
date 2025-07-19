package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/aggregator/client"
	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
	aggClient   client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, client client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka consumer: %w", err)
	}

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s: %w", topic, err)
	}
	return &KafkaConsumer{
		consumer:    c,
		calcService: svc,
		aggClient:   client,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka consumer started")
	c.isRunning = true
	c.consumeMessages()
}

func (c *KafkaConsumer) Stop() {
	c.isRunning = false
	c.consumer.Close()
	logrus.Info("Kafka consumer stopped")
}

func (c *KafkaConsumer) consumeMessages() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error: %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			continue
		}

		req := &types.AggregatorReq{
			ObuId: int64(data.OBUID),
			Value: c.calcService.CalculateDistance(&data),
			Unix:  time.Now().UnixNano(),
		}
		if err := c.aggClient.Aggregate(context.Background(), req); err != nil {
			logrus.WithFields(logrus.Fields{
				"obuID": data.OBUID,
				"error": err,
			}).Error("Failed to aggregate distance data")
			continue
		}
	}
}
