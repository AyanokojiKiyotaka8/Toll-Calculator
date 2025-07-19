package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"github.com/gorilla/websocket"
)

const (
	wsEndpoint       = "ws://localhost:30000/ws"
	numSimulatedOBUs = 10
)

var (
	sendInterval = time.Second
	r            *rand.Rand
)

func genCoord() float64 {
	n := float64(r.Intn(100) + 1)
	f := r.Float64()
	return n + f
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}

func genOBUIDs(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i] = r.Intn(math.MaxInt)
	}
	return ids
}

func main() {
	fmt.Println("running obu")
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}
	defer conn.Close()

	ids := genOBUIDs(numSimulatedOBUs)

	for {
		for _, id := range ids {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: id,
				Lat:   lat,
				Long:  long,
			}
			log.Printf("Sending OBU data: %+v", data)
			if err := conn.WriteJSON(data); err != nil {
				log.Printf("WebSocket write error: %v", err)
				continue
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}
