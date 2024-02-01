package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

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

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
