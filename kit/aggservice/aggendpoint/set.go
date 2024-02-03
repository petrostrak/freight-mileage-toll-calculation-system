package aggendpoint

import "github.com/go-kit/kit/endpoint"

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}
