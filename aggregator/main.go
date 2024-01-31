package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
)

func main() {
	httpAddr := flag.String("httpAddr", ":3010", "The listen address of the HTTP server")
	grpcAddr := flag.String("grpcAddr", ":3011", "The listen address of the gRPC server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)
	svc = NewMetricsMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcAddr, svc))
	}()

	log.Fatal(makeHTTPTransport(*httpAddr, svc))
}

func makeGRPCTransport(addr string, svc Aggregator) error {
	fmt.Println("gRPC transport running on port ", addr)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	server := grpc.NewServer([]grpc.ServerOption{}...)
	proto.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(addr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", addr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())

	return http.ListenAndServe(addr, nil)
}

func handleInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "method not supported"})
			return
		}

		obuID := r.URL.Query().Get("obuID")
		if obuID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBUID"})
			return
		}

		id, err := strconv.Atoi(obuID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBUID"})
			return
		}

		invoice, err := svc.CalculateInvoice(id)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "method not supported"})
			return
		}

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
