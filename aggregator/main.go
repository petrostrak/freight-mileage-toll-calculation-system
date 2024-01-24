package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

func main() {
	addr := flag.String("addr", ":3000", "The listen address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)

	if err := makeHTTPTransport(*addr, svc); err != nil {
		panic(err)
	}
}

func makeHTTPTransport(addr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", addr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice)

	return http.ListenAndServe(addr, nil)
}

func handleInvoice(w http.ResponseWriter, r *http.Request) {
	obuID := r.URL.Query().Get("obu_ID")
	if obuID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBUID"})
		return
	}
	w.Write([]byte(obuID))
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
