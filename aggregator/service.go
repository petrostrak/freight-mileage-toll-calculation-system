package main

import (
	"fmt"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
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
	fmt.Printf("processing and inserting distance in DB [ID:%d | Value: %.2f | Unix: %d]\n", distance.OBUID, distance.Value, distance.Unix)
	return i.store.Insert(distance)
}

func (i *InvoiceAggregator) CalculateInvoice(obu_ID int) (*types.Invoice, error) {
	distance, err := i.store.Get(obu_ID)
	if err != nil {
		return nil, err
	}

	invoice := &types.Invoice{
		OBUID:         obu_ID,
		TotalDistance: distance,
		TotalAmount:   baseFee * distance,
	}

	return invoice, nil
}
