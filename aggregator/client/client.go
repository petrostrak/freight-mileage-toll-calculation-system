package client

import (
	"context"

	"github.com/petrostrak/freight-mileage-toll-calculation-system/proto"
)

type Clienter interface {
	Aggregate(context.Context, *proto.AggregateRequest) error
}
