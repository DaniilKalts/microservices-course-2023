package interceptor

import (
	"context"
	"time"

	"github.com/sony/gobreaker/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreakerConfig struct {
	MaxRequests      uint32
	OpenTimeout      time.Duration
	FailureThreshold uint32
}

func CircuitBreakerInterceptor(cfg CircuitBreakerConfig) grpc.UnaryClientInterceptor {
	cb := gobreaker.NewTwoStepCircuitBreaker[any](gobreaker.Settings{
		Name:        "gateway-grpc",
		MaxRequests: cfg.MaxRequests,
		Timeout:     cfg.OpenTimeout,
		// When to open: after N failures in a row, the backend is likely down — stop sending traffic.
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= cfg.FailureThreshold
		},
		// What counts as failure: only real backend errors. If the client canceled or timed out,
		// that's not the backend's fault — don't count it toward tripping.
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			code := status.Code(err)
			return code == codes.Canceled || code == codes.DeadlineExceeded
		},
	})

	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		done, err := cb.Allow()
		if err != nil {
			return status.Error(codes.Unavailable, "circuit breaker is open")
		}

		invokeErr := invoker(ctx, method, req, reply, cc, opts...)
		done(invokeErr)

		return invokeErr
	}
}
