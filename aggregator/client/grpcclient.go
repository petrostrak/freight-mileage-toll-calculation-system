package client

import (
	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	Endpoint string
	client   proto.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	client := proto.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		client:   client,
	}, nil
}
