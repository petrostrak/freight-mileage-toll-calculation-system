package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetricHandler struct {
	reqCounter prometheus.Counter
}

func newHTTPMetricsHandler(reqName string) *HTTPMetricHandler {
	return &HTTPMetricHandler{
		reqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s_", reqName, "request_counter"),
			Name:      "aggregator",
		}),
	}
}

func handleInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
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
