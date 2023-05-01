package service

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

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

type MockExchangeService struct{}

func (MockExchangeService) GetConvertedTotal(target_currency string, amount int) (total float64, err error) {
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

func Test_Logging_GetConvertedTotal(t *testing.T) {
	tests := []struct {
		target_currency string
		amount          int
		msg             string
	}{
		{
			target_currency: "GBP",
			amount:          15,
			msg:             "method,GetConvertedTotal,target_currency,GBP,amount,15,total,2400,error,<nil>,duration",
		},
		{
			target_currency: "AAA",
			amount:          10,
			msg:             "method,GetConvertedTotal,target_currency,AAA,amount,10,total,0,error,Currency Not Found,duration",
		},
	}

	logger := new(MockLogger)
	var svc ExchangeService
	svc = new(MockExchangeService)
	svc = NewLoggingMiddleware(logger, svc)

	for id, test := range tests {
		svc.GetConvertedTotal(test.target_currency, test.amount)

		actual := logger.Result()

		assert.True(t, strings.HasPrefix(actual, test.msg), "~2|Test #%d logger expected: \"%s\", not: \"%s\"~", id, test.msg, actual)
	}
}
