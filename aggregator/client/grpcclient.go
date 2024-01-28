package client

import (
	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	Endpoint string
	proto.AggregatorClient
}

func NewGRPCClient(endpoint string) (*GRPCClient, error) {
	conn, err := grpc.Dial(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewAggregatorClient(conn)
	return &GRPCClient{
		Endpoint:         endpoint,
		AggregatorClient: client,
	}, nil
}
