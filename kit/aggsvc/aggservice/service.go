package aggservice

import (
	"context"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Service interface {
	Aggregate(context.Context, types.Distance) error
	Calculate(context.Context, int) (*types.Invoice, error)
}

const basePrice = 3.15

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
}

type BasicService struct {
	store Storer
}

func newBasicService(store Storer) Service {
	return &BasicService{store: store}
}

func (svc *BasicService) Aggregate(ctx context.Context, dist types.Distance) error {
	return svc.store.Insert(dist)
}

func (svc *BasicService) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	dist, err := svc.store.Get(obuID)
	if err != nil {
		return nil, err
	}

	inv := &types.Invoice{
		OBUID:         obuID,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}

	return inv, nil
}

// NewAggregatorService will construct a complete microservice
// with logging and instrumentation middleware.
func NewAggregatorService() Service {
	var svc Service
	{
		svc = newBasicService(NewMemoryStore())
		svc = newLoggingsMiddleware()(svc)
		svc = newInstrumentationMiddleware()(svc)
	}

	return svc
}
