package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

func main() {
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
		msg: make(chan types.OBUData, 128),
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
}

func (rcv *Receiver) wsReceiveLoop() {
	fmt.Println("New OBU connected client connected!")
	for {
		var data types.OBUData
		if err := rcv.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Printf("recv'd OBU data from [%d]:: <lat %.2f long %2.f>\n", data.OBUID, data.Lat, data.Long)
		rcv.msg <- data
	}
}
