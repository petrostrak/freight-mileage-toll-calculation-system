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

	return m.next.AggregateDistance(distance)
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
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})

	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregate_request_latency",
		Name:      "aggregate",
		Help:      "The request latency",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregate_request_latency",
		Name:      "calculate",
		Help:      "The request latency",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MetricsMiddleware{
		next:           next,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(distance types.Distance) error {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterAgg.Inc()
	}(time.Now())

	return m.next.AggregateDistance(distance)
}

func (m *MetricsMiddleware) CalculateInvoice(obuID int) (*types.Invoice, error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterCalc.Inc()
	}(time.Now())

	return m.next.CalculateInvoice(obuID)
}
