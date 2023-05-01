package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func LogTotalConvertedPriceEndpoint(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			logger.Log("endpoint", "TotalConvertedPriceEndpoint", "msg", "Calling endpoint")
			defer logger.Log("endpoint", "TotalConvertedPriceEndpoint", "msg", "Called endpoint")

			return next(ctx, request)
		}
	}
}
