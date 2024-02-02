package main

import (
	"context"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

type BasicService struct{}

func newBasicService() Service {
	return &BasicService{}
}

func (svc *BasicService) Aggregate(ctx context.Context, dist types.Distance) error {
	return nil
}

func (svc *BasicService) Calculate(context.Context, int) (*types.Invoice, error) {
	return nil, nil
}

// NewAggregatorService will construct a complete microservice
// with logging and instrumentation middleware.
func NewAggregatorService() Service {
	var svc Service
	{
		svc = newBasicService()
		svc = newLoggingsMiddleware()(svc)
		svc = newInstrumentationMiddleware()(svc)
	}

	return svc
}
