package main

import (
	"flag"
	"fmt"
	"go-app/repo"
	"go-app/service"
	"go-app/transport"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		listen = flag.String("listen", ":8080", "HTTP listen address")
	)
	flag.Parse()

	logger := log.NewLogfmtLogger(os.Stderr)

	fmt.Println("Repository: In progress")

	currencyRepo, _ := repo.NewCurrencyRepo("currency_rate.csv")

	fmt.Println("Repository: Ready")

	fmt.Println("Endpoints and handlers: In progress")

	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "gokitfundamentals",
		Subsystem: "pricing_service",
		Name:      "request_count",
		Help:      "Number of request received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "gokitfundamentals",
		Subsystem: "pricing_service",
		Name:      "request_latency",
		Help:      "Total duration of requests.",
	}, fieldKeys)

	var svc service.ExchangeService
	svc = service.NewExchangeService(currencyRepo)
	svc = service.NewInstrumentingMiddleware(requestCount, requestLatency, svc)
	svc = service.NewLoggingMiddleware(logger, svc)

	rtr := mux.NewRouter().StrictSlash(true)

	totalRetailPriceHandler := transport.MakeTotalConvertedPriceHttpHandler(logger, svc)
	rtr.Handle("/exchange", totalRetailPriceHandler).Methods(http.MethodPost)

	rtr.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)

	fmt.Println("Endpoints and handlers: Ready")

	fmt.Printf("Hosting on %s\n", *listen)

	http.ListenAndServe(*listen, rtr)
}
