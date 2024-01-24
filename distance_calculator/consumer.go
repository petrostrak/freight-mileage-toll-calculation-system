package main

import (
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/aggregator/client"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer          *kafka.Consumer
	isRunning         bool
	calculatorService Calculator
	aggregationClient *client.Client
}

func NewKafkaConsumer(topic string, srv Calculator, aggClient *client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:          c,
		calculatorService: srv,
		aggregationClient: aggClient,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error %s", err)
			continue
		}

		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			continue
		}

		distance, err := c.calculatorService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}

		req := types.Distance{
			Value: distance,
			OBUID: data.OBUID,
			Unix:  time.Now().UnixNano(),
		}

		if err = c.aggregationClient.AggregateInvoice(req); err != nil {
			logrus.Errorf("aggregate error: %s", err)
			continue
		}
	}
}
