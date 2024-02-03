package aggservice

import (
	"context"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
}

func newLoggingsMiddleware() Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
		}
	}
}

func (l loggingMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (l loggingMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}

type instrumentationMiddleware struct {
	next Service
}

func newInstrumentationMiddleware() Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{
			next: next,
		}
	}
}

func (i instrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (i instrumentationMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}
