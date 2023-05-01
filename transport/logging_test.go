package transport

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

type MockLogger struct {
	result []string
}

func (ml *MockLogger) Log(keyvals ...interface{}) (err error) {
	var result []string
	for _, val := range keyvals {
		result = append(result, fmt.Sprint(val))
	}

	ml.result = result

	return nil
}

func (ml *MockLogger) Result() string {
	return strings.Join(ml.result[:], ",")
}

func Test_LogTotalConvertedPriceEndpoint(t *testing.T) {
	tests := []struct {
		service  string
		request  ConvertedPriceRequest
		expected string
	}{
		{
			service:  "endpointConvertedTest",
			request:  ConvertedPriceRequest{"USD", 10},
			expected: "service,endpointConvertedTest,endpoint,TotalConvertedPriceEndpoint,msg,Called endpoint",
		},
		{
			service:  "endpointConvertedTest",
			request:  ConvertedPriceRequest{"GBP", 20},
			expected: "service,endpointConvertedTest,endpoint,TotalConvertedPriceEndpoint,msg,Called endpoint",
		},
	}

	endpoint := func(_ context.Context, request interface{}) (interface{}, error) {
		return request, nil
	}

	logger := &MockLogger{}

	for id, test := range tests {
		lmw := LogTotalConvertedPriceEndpoint(log.With(logger, "service", test.service))(endpoint)

		lmw(context.Background(), test.request)

		actual := logger.Result()

		assert.True(t, test.expected == actual, "~2|Test #%d logger expected: \"%s\", not: \"%s\"~", id, test.expected, actual)
	}
}
