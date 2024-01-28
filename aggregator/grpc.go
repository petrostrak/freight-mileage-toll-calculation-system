package main

import (
	"context"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
)

type GRPCAggregatorServer struct {
	proto.UnimplementedAggregatorServer
	svc Aggregator
}

func NewGRPCAggregatorServer(svc Aggregator) *GRPCAggregatorServer {
	return &GRPCAggregatorServer{
		svc: svc,
	}
}

func (s *GRPCAggregatorServer) Aggregate(ctx context.Context, req *proto.AggregateRequest) (*proto.None, error) {
	distance := types.Distance{
		OBUID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &proto.None{}, s.svc.AggregateDistance(distance)
}
