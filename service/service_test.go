package service

import (
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockProductRepo struct{}

func (MockProductRepo) FetchPrice(target_currency string) (price float64, found bool) {
	data := []string{
		"GBP,160",
		"USD,120",
		"EUR,140",
	}

	for _, line := range data {
		parts := strings.Split(line, ",")
		if parts[0] != target_currency {
			continue
		}

		price, _ = strconv.ParseFloat(parts[1], 64)

		return price, true
	}

	return 0, false
}

func Test_GetConvertedTotal(t *testing.T) {
	tests := []struct {
		target_currency string
		amount          int
		err             error
		total           float64
	}{
		{
			target_currency: "",
			amount:          0,
			err:             ErrInvalidTargetCurrency,
			total:           0.0,
		},
		{
			target_currency: "GBP",
			amount:          0,
			err:             ErrInvalidAmount,
			total:           0.0,
		},
		{
			target_currency: "USD",
			amount:          15,
			err:             nil,
			total:           1800.00,
		},
		{
			target_currency: "AAA",
			amount:          10,
			err:             ErrTargetCurrencyNotFound,
			total:           0.0,
		},
	}

	mockProductRepo := new(MockProductRepo)

	convertService := NewExchangeService(mockProductRepo)

	for id, test := range tests {
		total, err := convertService.GetConvertedTotal(test.target_currency, test.amount)
		assert.True(t, test.err == err, "~2|Test #%d expected error: %s, not error %s~", id, test.err, err)
		assert.True(t, test.total == total, "~2|Test #%d expected total: %.2f, not total %.2f~", id, test.total, total)
	}
}
