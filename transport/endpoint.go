package transport

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type ExchangeService interface {
	GetConvertedTotal(target_currency string, amount int) (total float64, err error)
}

func MakeTotalConvertedPriceEndpoint(svc ExchangeService) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		req := request.(ConvertedPriceRequest)
		total, err := svc.GetConvertedTotal(req.TargetCurrency, req.Amount)
		if err != nil {
			return ConvertedPriceResponse{0.0, err.Error()}, nil
		}

		return ConvertedPriceResponse{total, ""}, nil
	}
}
