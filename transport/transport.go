package transport

import (
	"context"
	"encoding/json"
	"net/http"

	gkendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
)

const (
	INVALID_REQUEST = "Invalid Request"
)

func decodeConvertedPriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request ConvertedPriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, &ErrorResponse{Err: INVALID_REQUEST}
	}

	return request, nil
}

func MakeTotalConvertedPriceHttpHandler(logger log.Logger, svc ExchangeService) *httptransport.Server {
	var retailEndpoint gkendpoint.Endpoint
	retailEndpoint = MakeTotalConvertedPriceEndpoint(svc)
	retailEndpoint = LogTotalConvertedPriceEndpoint(log.With(logger, "service", "ExchangeService"))(retailEndpoint)

	return httptransport.NewServer(
		retailEndpoint,
		decodeConvertedPriceRequest,
		encodeResponse,
	)
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
