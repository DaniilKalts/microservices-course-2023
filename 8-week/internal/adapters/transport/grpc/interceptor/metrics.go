package interceptor

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	requestStatusSuccess = "success"
	requestStatusError   = "error"
)

func MetricsInterceptor(registry prometheus.Registerer) grpc.UnaryServerInterceptor {
	factory := promauto.With(registry)

	responseCounter := factory.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "termachat",
			Subsystem: "grpc",
			Name:      "responses_total",
			Help:      "Number of gRPC responses by method and status",
		},
		[]string{"method", "status"},
	)

	requestDuration := factory.NewHistogramVec(
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

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		resp, err := handler(ctx, req)

		reqStatus := responseStatus(err)
		method := info.FullMethod

		responseCounter.WithLabelValues(method, reqStatus).Inc()
		requestDuration.WithLabelValues(method, reqStatus).Observe(time.Since(startedAt).Seconds())

		return resp, err
	}
}

func responseStatus(err error) string {
	if status.Code(err) == codes.OK {
		return requestStatusSuccess
	}

	return requestStatusError
}
