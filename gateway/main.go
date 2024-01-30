package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/aggregator/client"
	"github.com/sirupsen/logrus"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	addr := flag.String("addr", ":6001", "the listen address of the gateway HTTP server")
	aggregatorServiceAddr := flag.String("aggServiceAddr", "http://127.0.0.1:3000", "the address of the aggregator service")
	flag.Parse()

	client := client.NewHTTPClient(*aggregatorServiceAddr)
	invoiceHandler := newInvoiceHandler(client)

	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.handleGetInvoice))
	logrus.Infof("gateway HTTP server running on port %s", *addr)
	log.Fatal(http.ListenAndServe(":6001", nil))
}

type InvoiceHandler struct {
	client client.Clienter
}

func newInvoiceHandler(client client.Clienter) *InvoiceHandler {
	return &InvoiceHandler{client}
}

func (h *InvoiceHandler) handleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	// access aggregator client
	inv, err := h.client.GetInvoice(context.Background(), 1243124)
	if err != nil {
		return err
	}

	return writeJSON(w, http.StatusOK, map[string]any{
		"invoice": inv,
	})
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("req::")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]any{
				"error": err,
			})
		}
	}
}
