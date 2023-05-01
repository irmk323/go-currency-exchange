package transport

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
)

var (
	ErrInvalidTargetCurrency  = errors.New("Invalid Target Currency Requested")
	ErrTargetCurrencyNotFound = errors.New("TargetCurrency Not Found")
	ErrInvalidAmount          = errors.New("Invalid Amount Requested")
)

type MockPricingService struct{}

func (MockPricingService) GetConvertedTotal(target_currency string, amount int) (total float64, err error) {
	if target_currency == "" {
		return 0.0, ErrInvalidTargetCurrency
	}
	if amount <= 0 {
		return 0.0, ErrInvalidAmount
	}

	data := []string{
		"GBP,160",
		"USD,120",
		"EUR,140",
	}

	for _, line := range data {
		parts := strings.Split(line, ",")
		if parts[0] == target_currency {
			price, _ := strconv.ParseFloat(parts[1], 64)

			return (price * float64(amount)), nil
		}
	}

	return 0.0, ErrTargetCurrencyNotFound
}

func Test_MakeConvertedPriceEndpoint(t *testing.T) {
	tests := []struct {
		request  ConvertedPriceRequest
		response ConvertedPriceResponse
	}{
		{
			request:  ConvertedPriceRequest{TargetCurrency: "", Amount: 0},
			response: ConvertedPriceResponse{Err: "Invalid Target Currency Requested"},
		},
		{
			request:  ConvertedPriceRequest{TargetCurrency: "GBP", Amount: 0},
			response: ConvertedPriceResponse{Err: "Invalid Amount Requested"},
		},
		{
			request:  ConvertedPriceRequest{TargetCurrency: "USD", Amount: 15},
			response: ConvertedPriceResponse{Total: 1800.00},
		},
		{
			request:  ConvertedPriceRequest{TargetCurrency: "AAA", Amount: 10},
			response: ConvertedPriceResponse{Err: "TargetCurrency Not Found"},
		},
	}

	mockPricingService := new(MockPricingService)

	totalConvertedPriceHandler := httptransport.NewServer(
		MakeTotalConvertedPriceEndpoint(mockPricingService),
		decodeConvertedPriceRequest,
		encodeResponse,
	)

	server := httptest.NewServer(totalConvertedPriceHandler)
	defer server.Close()

	for id, test := range tests {
		postBody, _ := json.Marshal(test.request)

		responseBody := bytes.NewBuffer(postBody)
		resp, err := http.Post(server.URL, "application/json", responseBody)
		if err != nil {
			log.Fatalf("An Error Occured %v", err)
		}

		var actualResponse ConvertedPriceResponse
		json.NewDecoder(resp.Body).Decode(&actualResponse)

		assert.True(t, test.response.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, test.response.Err, actualResponse.Err)
		assert.True(t, test.response.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.response.Total, actualResponse.Total)
	}
}
