package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var ResponseCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Namespace: "termachat",
		Subsystem: "grpc",
		Name:      "responses_total",
		Help:      "Number of gRPC responses by method and status",
	},
	[]string{"method", "status"},
)

var RequestDuration = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Namespace: "termachat",
		Subsystem: "grpc",
		Name:      "duration_seconds",
		Help:      "gRPC request duration in seconds",
		Buckets: []float64{
			0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5,
		},
	},
	[]string{"method", "status"},
)

func init() {
	prometheus.MustRegister(ResponseCounter, RequestDuration)
}
