package transport

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_TotalConvertedPriceRequest(t *testing.T) {
	tests := []struct {
		input    ConvertedPriceRequest
		expected ConvertedPriceRequest
	}{
		{
			input:    ConvertedPriceRequest{TargetCurrency: "test", Amount: 0},
			expected: ConvertedPriceRequest{TargetCurrency: "test", Amount: 0},
		},
		{
			input:    ConvertedPriceRequest{TargetCurrency: "", Amount: 12},
			expected: ConvertedPriceRequest{TargetCurrency: "", Amount: 12},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual ConvertedPriceRequest
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.TargetCurrency == actual.TargetCurrency, "~2|Test #%d expected code: %s, not code %s~", id, test.expected.TargetCurrency, actual.TargetCurrency)
		assert.True(t, test.expected.Amount == actual.Amount, "~2|Test #%d expected qty: %d, not qty %d~", id, test.expected.Amount, actual.Amount)
	}
}

func Test_TotalRetailPriceResponse(t *testing.T) {
	tests := []struct {
		input    ConvertedPriceResponse
		expected ConvertedPriceResponse
	}{
		{
			input:    ConvertedPriceResponse{Total: 100.99},
			expected: ConvertedPriceResponse{Total: 100.99},
		},
		{
			input:    ConvertedPriceResponse{Total: 0.0, Err: "test"},
			expected: ConvertedPriceResponse{Total: 0.0, Err: "test"},
		},
	}

	for id, test := range tests {
		data, _ := json.Marshal(test.input)

		var actual ConvertedPriceResponse
		json.Unmarshal(data, &actual)

		assert.True(t, test.expected.Total == actual.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.expected.Total, actual.Total)
		assert.True(t, test.expected.Err == actual.Err, "~2|Test #%d expected err: %s, not err %s~", id, test.expected.Err, actual.Err)
	}
}

func Test_MakeTotalConvertedPriceHttpHandler(t *testing.T) {
	tests := []struct {
		request  interface{}
		response interface{}
	}{
		{
			request:  ConvertedPriceRequest{TargetCurrency: "", Amount: 0},
			response: ConvertedPriceResponse{Err: "Invalid Target Currency Requested"},
		},
		{
			request:  ConvertedPriceRequest{TargetCurrency: "USD", Amount: 0},
			response: ConvertedPriceResponse{Err: "Invalid Amount Requested"},
		},
		{
			request:  ConvertedPriceRequest{TargetCurrency: "GBP", Amount: 15},
			response: ConvertedPriceResponse{Total: 2400.00},
		},
		{
			request:  ConvertedPriceRequest{TargetCurrency: "fff000", Amount: 10},
			response: ConvertedPriceResponse{Err: "TargetCurrency Not Found"},
		},
		{
			request:  "test",
			response: ConvertedPriceResponse{Err: "Invalid Request"},
		},
	}

	mockPricingService := new(MockPricingService)

	logger := &MockLogger{}
	totalConvertedPriceHandler := MakeTotalConvertedPriceHttpHandler(logger, mockPricingService)

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

		testResponse := test.response.(ConvertedPriceResponse)

		assert.True(t, testResponse.Err == actualResponse.Err, "~2|Test #%d expected error: %s, not error %s~", id, testResponse.Err, actualResponse.Err)
		assert.True(t, testResponse.Total == actualResponse.Total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, testResponse.Total, actualResponse.Total)
	}
}
