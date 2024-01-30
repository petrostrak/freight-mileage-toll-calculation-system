package main

import (
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/sirupsen/logrus"
)

const baseFee = 3.45

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	logrus.WithFields(logrus.Fields{
		"obuID":     distance.OBUID,
		"obu_Value": distance.Value,
		"obu_Unix":  distance.Unix,
	}).Info("processing and inserting distance in DB")
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obuID int) (*types.Invoice, error) {
	distance, err := i.store.Get(obuID)
	if err != nil {
		return nil, err
	}

	invoice := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: distance,
		TotalAmount:   baseFee * distance,
	}

	return invoice, nil
}
