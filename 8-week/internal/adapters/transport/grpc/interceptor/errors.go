package interceptor

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorRule struct {
	Target  error
	Code    codes.Code
	Message string
}

func ErrorMappingInterceptor(logger *zap.Logger, rules []ErrorRule) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err == nil {
			return resp, nil
		}

		if _, ok := status.FromError(err); ok {
			return resp, err
		}

		for _, r := range rules {
			if errors.Is(err, r.Target) {
				msg := r.Message
				if msg == "" {
					msg = r.Target.Error()
				}
				return nil, status.Error(r.Code, msg)
			}
		}

		logger.Error("unmapped error", zap.String("method", info.FullMethod), zap.Error(err))
		return nil, status.Error(codes.Internal, "internal error")
	}
}
