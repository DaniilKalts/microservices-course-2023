package interceptor

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validatable interface {
	Validate() error
}

func ValidationInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		if reqWithValidation, ok := req.(validatable); ok {
			if err := reqWithValidation.Validate(); err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "%s: %v", info.FullMethod, err)
			}
		}

		return handler(ctx, req)
	}
}
