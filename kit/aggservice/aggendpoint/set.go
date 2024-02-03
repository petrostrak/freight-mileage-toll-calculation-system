package aggendpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

type AggregateRequest struct {
	OBUID int     `json:"obuID"`
	Value float64 `json:"value"`
	Unix  int64   `json:"unix"`
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})
	if err != nil {
		return err
	}

	return nil
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

func (s Set) Calculate(ctx context.Context, obuID int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{
		OBUID: obuID,
	})
	if err != nil {
		return nil, err
	}

	result := resp.(CalculateResponse)

	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}
