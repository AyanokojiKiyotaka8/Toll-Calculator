package main

import (
	"log"
	"net/http"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/gorilla/websocket"
)

func main() {
	receiver, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", receiver.handleWs)
	log.Println("WebSocket server started on :30000")
	if err := http.ListenAndServe(":30000", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}

type DataReceiver struct {
	prod DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p          DataProducer
		err        error
		kafkaTopic = "obudata"
	)
	p, err = NewKafkaDataProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	return &DataReceiver{
		prod: p,
	}, nil
}

func (dr *DataReceiver) handleWs(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		http.Error(w, "WebSocket upgrade failed", http.StatusBadRequest)
		return
	}

	go dr.wsReceiveLoop(conn)
}

func (dr *DataReceiver) wsReceiveLoop(conn *websocket.Conn) {
	defer conn.Close()
	defer dr.prod.Stop()

	for {
		var data types.OBUData
		if err := conn.ReadJSON(&data); err != nil {
			log.Println("Client disconnected or read error:", err)
			continue
		}

		if err := dr.prod.ProduceData(&data); err != nil {
			log.Println("Kafka produce error:", err)
		}
	}
}
