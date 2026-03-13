package diagnostic

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HealthChecker interface {
	Name() string
	Check(ctx context.Context) error
}

const (
	readTimeout  = 5 * time.Second
	writeTimeout = 5 * time.Second
)

type Deps struct {
	Address  string
	Checkers []HealthChecker
	Gatherer prometheus.Gatherer
}

func NewServer(deps Deps) *http.Server {
	mux := http.NewServeMux()

	metricsHandler := promhttp.Handler()
	if deps.Gatherer != nil {
		metricsHandler = promhttp.HandlerFor(deps.Gatherer, promhttp.HandlerOpts{})
	}
	mux.Handle("/metrics", metricsHandler)

	mux.HandleFunc("/healthz/liveness", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	mux.HandleFunc("/healthz/readiness", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		for _, c := range deps.Checkers {
			if err := c.Check(ctx); err != nil {
				w.WriteHeader(http.StatusServiceUnavailable)
				_, _ = w.Write([]byte(fmt.Sprintf("%s: %s", c.Name(), err.Error())))
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	return &http.Server{
		Addr:         deps.Address,
		Handler:      mux,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}
