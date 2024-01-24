package main

import (
	"time"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
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

func (m *LogMiddleware) CalculateInvoice(obu_id int) (invoice *types.Invoice, err error) {
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
			"obu_ID":         obu_id,
			"total_distance": distance,
			"total_amount":   amount,
		}).Info("Calculate invoice")
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(obu_id)
	return
}
