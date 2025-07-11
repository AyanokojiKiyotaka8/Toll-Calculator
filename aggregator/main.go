package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AyanokojiKiyotaka8/Toll-Calculator/types"
)

func main() {
	var svc Aggregator
	store := NewMemoryStore()
	svc = NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	makeHTTPTransport(svc, ":3000")
}

func makeHTTPTransport(svc Aggregator, listenAddress string) {
	fmt.Println("starting HTTP Transport")
	http.HandleFunc("/aggregate", handleAggregateDistance(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.ListenAndServe(listenAddress, nil)
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
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
