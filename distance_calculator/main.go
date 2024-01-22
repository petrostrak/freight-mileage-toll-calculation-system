package main

import (
	"log"
)

// type DistanceCalculator struct {
// 	consumer DataConsumer
// }

const kafkaTopic = "obuData"

func main() {
	kafkaConsmuser, err := NewKafkaConsumer(kafkaTopic)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsmuser.Start()
}
