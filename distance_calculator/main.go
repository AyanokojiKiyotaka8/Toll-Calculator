package main

import (
	"log"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/aggregator/client"
)

const (
	kafkaTopic   = "obudata"
	httpEndPoint = "http://127.0.0.1:3000/aggregate"
	grpcEndPoint = "127.0.0.1:3001"
)

func main() {
	var svc CalculatorServicer
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	// TODO: Replace gRPC client with HTTP client if needed
	// httpClient := client.NewHTTPClient(httpEndPoint)
	grpcClient, err := client.NewGRPCClient(grpcEndPoint)
	if err != nil {
		log.Fatal(err)
	}
	defer grpcClient.Close()

	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, grpcClient)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
	kafkaConsumer.Stop()
}
