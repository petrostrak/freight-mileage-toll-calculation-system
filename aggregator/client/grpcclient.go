package client

import (
	"context"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	client   proto.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint: endpoint,
		client:   client,
	}, nil
}

func (c *GRPCClient) Aggregate(ctx context.Context, req *proto.AggregateRequest) error {
	_, err := c.client.Aggregate(ctx, req)
	return err
}
