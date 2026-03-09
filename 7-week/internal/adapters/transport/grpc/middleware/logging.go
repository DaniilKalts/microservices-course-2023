package middleware

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		startedAt := time.Now()

		resp, err := handler(ctx, req)
		duration := time.Since(startedAt)

		code := status.Code(err)

		remoteAddr := ""
		if p, ok := peer.FromContext(ctx); ok && p != nil && p.Addr != nil {
			remoteAddr = p.Addr.String()
		}

		fields := []zap.Field{
			zap.String("protocol", "grpc"),
			zap.String("method", info.FullMethod),
			zap.Int("status_code", int(code)),
			zap.String("remote_addr", remoteAddr),
			zap.NamedError("error", err),
			zap.Float64("duration_ms", float64(duration)/float64(time.Millisecond)),
		}

		switch code {
		case codes.OK:
			logger.Info("request completed", fields...)
		case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
			logger.Error("request completed", fields...)
		default:
			logger.Warn("request completed", fields...)
		}

		return resp, err
	}
}
