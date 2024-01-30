package main

import (
	"time"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) Aggregator {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("aggregate distance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleware) CalculateInvoice(obuID int) (invoice *types.Invoice, err error) {
	var (
		distance float64
		amount   float64
	)

	if invoice != nil {
		distance = invoice.TotalDistance
		amount = invoice.TotalAmount
	}

	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":           time.Since(start),
			"err":            err,
			"obuID":          obuID,
			"total_distance": distance,
			"total_amount":   amount,
		}).Info("Calculate invoice")
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(obuID)
	return
}

type MetricsMiddleware struct {
	reqCounter prometheus.Counter
	reqLatency prometheus.Histogram
	next       Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	reqCounter := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregator_request_counter",
	})
	reqLatency := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregator_request_counter",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:       next,
		reqCounter: reqCounter,
		reqLatency: reqLatency,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) error {
	return m.next.AggregateDistance(distance)
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (*types.Invoice, error) {
	return m.next.CalculateInvoice(obuID)
}
