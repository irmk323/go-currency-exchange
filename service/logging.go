package service

import (
	"time"

	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   ExchangeService
}

func NewLoggingMiddleware(logger log.Logger, next ExchangeService) (lmw *loggingMiddleware) {
	lmw = &loggingMiddleware{
		logger: logger,
		next:   next,
	}

	return
}

func (mw loggingMiddleware) GetConvertedTotal(target_currency string, amount int) (total float64, err error) {
	defer func(begin time.Time) {
		_ = mw.logger.Log(
			"method", "GetConvertedTotal",
			"target_currency", target_currency,
			"amount", amount,
			"total", total,
			"error", err,
			"duration", time.Since(begin),
		)
	}(time.Now())

	total, err = mw.next.GetConvertedTotal(target_currency, amount)

	return
}
