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

func (l loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return l.next.Aggregate(ctx, dist)
}

func (l loggingMiddleware) Calculate(ctx context.Context, dist int) (*types.Invoice, error) {
	return l.next.Calculate(ctx, dist)
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

func (i instrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return i.next.Aggregate(ctx, dist)
}

func (i instrumentationMiddleware) Calculate(ctx context.Context, dist int) (*types.Invoice, error) {
	return i.next.Calculate(ctx, dist)
}
