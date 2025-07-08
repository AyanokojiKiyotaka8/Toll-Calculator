package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Println("running receiver")
	dr := NewDataReceiver()
	http.HandleFunc("/ws", dr.handleWs)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	conn *websocket.Conn
}

func NewDataReceiver() *DataReceiver {
	return &DataReceiver{}
}

func (dr *DataReceiver) handleWs(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn
	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error: ", err)
			continue
		}
		fmt.Println(data)
	}
}
