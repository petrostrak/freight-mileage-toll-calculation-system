package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/obu/types"
	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
)

type HTTPClient struct {
	Endpoint string
}

func NewHTTPClient(endpoint string) *HTTPClient {
	return &HTTPClient{
		Endpoint: endpoint,
	}
}

func (c *HTTPClient) GetInvoice(ctx context.Context, id int) (*types.Invoice, error) {
	return &types.Invoice{
		OBUID:         id,
		TotalDistance: 505.6,
		TotalAmount:   41.000,
	}, nil
}

func (c *HTTPClient) Aggregate(ctx context.Context, aggregationReq *proto.AggregateRequest) error {
	b, err := json.Marshal(aggregationReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the service responded with non 200 status code %d", resp.StatusCode)
	}

	return nil
}
