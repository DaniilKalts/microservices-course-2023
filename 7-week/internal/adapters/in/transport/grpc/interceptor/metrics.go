package interceptor

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	requestStatusSuccess = "success"
	requestStatusError   = "error"
)

func MetricsInterceptor(responseCounter *prometheus.CounterVec, requestDuration *prometheus.HistogramVec) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		resp, err := handler(ctx, req)

		requestStatus := responseStatus(err)
		method := info.FullMethod

		responseCounter.WithLabelValues(method, requestStatus).Inc()
		requestDuration.WithLabelValues(method, requestStatus).Observe(time.Since(startedAt).Seconds())

		return resp, err
	}
}

func responseStatus(err error) string {
	if status.Code(err) == codes.OK {
		return requestStatusSuccess
	}

	return requestStatusError
}
