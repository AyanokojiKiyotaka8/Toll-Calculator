package main

import (
	"encoding/json"
	"fmt"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

type KafKaConsumer struct {
	consumer    *kafka.Consumer
	isRunning   bool
	calcService CalculatorServicer
}

func NewKafkaConsumer(topic string, svc CalculatorServicer) (*KafKaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	if err := c.SubscribeTopics([]string{topic}, nil); err != nil {
		return nil, err
	}
	return &KafKaConsumer{
		consumer:    c,
		calcService: svc,
	}, nil
}

func (c *KafKaConsumer) Start() {
	logrus.Info("kafka consumer started")
	c.isRunning = true
	c.consumeMessages()
}

func (c *KafKaConsumer) consumeMessages() {
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

		distance := c.calcService.CalculateDistance(&data)
		fmt.Println(distance)
	}
}
