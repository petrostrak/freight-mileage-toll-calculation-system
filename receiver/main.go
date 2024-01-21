package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gorilla/websocket"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	for i := 0; i < 10; i++ {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte("test producing message"),
		}, nil)
	}

	recv := NewReceiver()
	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":30000", nil)
}

type Receiver struct {
	msg  chan types.OBUData
	conn *websocket.Conn
}

func NewReceiver() *Receiver {
	return &Receiver{
		msg: make(chan types.OBUData),
	}
}

func (rcv *Receiver) wsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	rcv.conn = conn
	go rcv.wsReceiveLoop()
	for data := range rcv.msg {
		fmt.Printf("recv'd OBU data from [%d]:: <lat %.2f long %2.f>\n", data.OBUID, data.Lat, data.Long)
	}
}

func (rcv *Receiver) wsReceiveLoop() {
	fmt.Println("New OBU connected client connected!")
	for {
		var data types.OBUData
		if err := rcv.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		rcv.msg <- data
	}
}
