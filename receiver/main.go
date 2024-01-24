package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

func main() {
	recv, err := NewReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", recv.wsHandler)
	http.ListenAndServe(":30000", nil)
}

type Receiver struct {
	// msg  chan types.OBUData
	conn *websocket.Conn
	prod DataProducer
}

func NewReceiver() (*Receiver, error) {
	var kafkaTopic = "obuData"
	p, err := NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}

	return &Receiver{
		// msg:  make(chan types.OBUData),
		prod: NewLogMiddleware(p),
	}, nil
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
}

func (rcv *Receiver) wsReceiveLoop() {
	fmt.Println("New OBU client connected!")
	for {
		var data types.OBUData
		if err := rcv.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		if err := rcv.produceData(data); err != nil {
			fmt.Printf("kafka produceData error: %s\n", err)
		}
	}
}

func (rcv *Receiver) produceData(data types.OBUData) error {
	return rcv.prod.ProduceData(data)
}
