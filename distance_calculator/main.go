package main

import (
	"log"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const kafkaTopic = "obuData"

func main() {
	srv := NewCalculatorService()
	srv = NewLogMiddleware(srv)
	kafkaConsmuser, err := NewKafkaConsumer(kafkaTopic, srv)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsmuser.Start()
}
