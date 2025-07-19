package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
	"google.golang.org/grpc"
)

func main() {
	var svc Aggregator
	store := NewMemoryStore()
	svc = NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)

	go makeGRPCTransport(svc, ":3001")
	makeHTTPTransport(svc, ":3000")
}

func makeGRPCTransport(svc Aggregator, listenAddress string) {
	fmt.Println("starting GRPC Transport on port", listenAddress)
	aggServer := NewGRPCAggregatorServer(svc)
	grpcServer := grpc.NewServer()
	types.RegisterAggregatorServer(grpcServer, aggServer)

	l, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	err = grpcServer.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}

func makeHTTPTransport(svc Aggregator, listenAddress string) {
	fmt.Println("starting HTTP Transport on port", listenAddress)
	http.HandleFunc("/aggregate", handleAggregateDistance(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Fatal(err)
	}
}

func handleAggregateDistance(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vals, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu ID"})
			return
		}
		obuID, err := strconv.Atoi(vals[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid obu ID"})
			return
		}
		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"invoice": invoice})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
