package auth

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

func AuthInterceptor(jwtManager jwt.Manager, policy AccessPolicy) grpc.UnaryServerInterceptor {
	if policy.IsEmpty() {
		return func(
			_ context.Context,
			_ any,
			_ *grpc.UnaryServerInfo,
			_ grpc.UnaryHandler,
		) (any, error) {
			return nil, status.Error(codes.Internal, domainAuth.ErrAccessPolicyNotConfigured.Error())
		}
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requiredGroup, ok := policy.GroupForMethod(info.FullMethod)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, domainAuth.ErrAccessDenied.Error())
		}

		if requiredGroup == AccessGroupPublic {
			return handler(ctx, req)
		}

		token, err := accessTokenFromContext(ctx)
		if err != nil {
			return nil, err
		}

		claims, err := authorize(token, jwtManager, requiredGroup)
		if err != nil {
			return nil, err
		}

		ctx = withClaims(ctx, claims)

		return handler(ctx, req)
	}
}
