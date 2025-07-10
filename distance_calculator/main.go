package main

import "log"

const kafkaTopic = "obudata"

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
