package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/metrics"
)

const (
	requestStatusSuccess = "success"
	requestStatusError   = "error"
)

func MetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()
		metrics.RequestCounter.Inc()

		resp, err := handler(ctx, req)

		requestStatus := responseStatus(err)
		method := info.FullMethod

		metrics.ResponseCounter.WithLabelValues(method, requestStatus).Inc()
		metrics.RequestDuration.WithLabelValues(method, requestStatus).Observe(time.Since(startedAt).Seconds())

		return resp, err
	}
}

func responseStatus(err error) string {
	if status.Code(err) == codes.OK {
		return requestStatusSuccess
	}

	return requestStatusError
}
