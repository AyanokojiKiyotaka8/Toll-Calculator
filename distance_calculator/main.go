package main

import (
	"log"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/aggregator/client"
)

const (
	kafkaTopic = "obudata"
	endPoint   = "http://127.0.0.1:3000/aggregate"
)

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.NewClient(endPoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
