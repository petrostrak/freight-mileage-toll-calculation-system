package aggendpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

func LoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request any) (response any, err error) {
			defer func(start time.Time) {
				logger.Log("error", err, "took", time.Since(start))
			}(time.Now())
			return next(ctx, request)
		}
	}
}
