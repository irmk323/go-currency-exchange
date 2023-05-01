package service

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           ExchangeService
}

func NewInstrumentingMiddleware(requestCount metrics.Counter, requestLatency metrics.Histogram, next ExchangeService) (imw *instrumentingMiddleware) {
	imw = &instrumentingMiddleware{
		requestCount:   requestCount,
		requestLatency: requestLatency,
		next:           next,
	}

	return
}

func (mw instrumentingMiddleware) GetConvertedTotal(code string, qty int) (total float64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetConvertedTotal", "error", fmt.Sprint(err != nil)}

		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	total, err = mw.next.GetConvertedTotal(code, qty)

	return
}
