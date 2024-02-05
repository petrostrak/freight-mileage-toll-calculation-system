package aggservice

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	log  log.Logger
	next Service
}

func newLoggingsMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			log:  logger,
			next: next,
		}
	}
}

func (l loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		l.log.Log("took:", time.Since(start), "obuID", dist.OBUID, "distance", dist.Value, "err", err)
	}(time.Now())
	err = l.next.Aggregate(ctx, dist)
	return
}

func (l loggingMiddleware) Calculate(ctx context.Context, dist int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		l.log.Log("took:", time.Since(start), "inv", inv, "err", err)
	}(time.Now())
	inv, err = l.next.Calculate(ctx, dist)
	return
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
