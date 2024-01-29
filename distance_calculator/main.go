package main

import (
	"log"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/aggregator/client"
)

const (
	kafkaTopic         = "obuData"
	aggregatorEndpoint = "http://127.0.0.1:3000"
)

func main() {
	srv := NewCalculatorService()
	srv = NewLogMiddleware(srv)
	kafkaConsmuser, err := NewKafkaConsumer(kafkaTopic, srv, client.NewHTTPClient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsmuser.Start()
}
