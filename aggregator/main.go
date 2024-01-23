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

	if err := makeHTTPTransport(*addr, svc); err != nil {
		panic(err)
	}
}

func makeHTTPTransport(addr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", addr)
	http.HandleFunc("/aggregate", handleAggregate(svc))

	return http.ListenAndServe(addr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
