package auth

import (
	"context"
	"errors"
	"strings"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
)

var (
	ErrAuthorizationTokenMissing = errors.New("authorization token is required")
	ErrInvalidAccessToken        = errors.New("invalid access token")
	ErrInsufficientPermissions   = errors.New("insufficient role permissions")
	ErrMethodNotRegistered       = errors.New("method not registered in access policy")
)

func NewInterceptor(jwtManager jwt.Manager, policy AccessPolicy, logger *zap.Logger) (grpc.UnaryServerInterceptor, error) {
	if jwtManager == nil {
		return nil, errors.New("jwt manager is nil")
	}
	if policy.IsEmpty() {
		return nil, errors.New("access policy is empty")
	}

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		group, ok := policy.GroupForMethod(info.FullMethod)
		if !ok {
			logger.Warn("method not registered in access policy", zap.String("method", info.FullMethod))
			return nil, status.Error(codes.PermissionDenied, ErrMethodNotRegistered.Error())
		}

		if group.IsPublic {
			return handler(ctx, req)
		}

		token, err := accessTokenFromContext(ctx)
		if err != nil {
			logger.Warn("missing authorization token", zap.String("method", info.FullMethod))
			return nil, err
		}

		claims, err := jwtManager.VerifyAccessToken(token)
		if err != nil || claims == nil {
			logger.Warn("authorization failed",
				zap.String("method", info.FullMethod),
				zap.Error(err),
			)
			return nil, status.Error(codes.Unauthenticated, ErrInvalidAccessToken.Error())
		}

		if !group.AllowsRole(claims.RoleID) {
			logger.Warn("authorization failed",
				zap.String("method", info.FullMethod),
				zap.Int32("role_id", claims.RoleID),
			)
			return nil, status.Error(codes.PermissionDenied, ErrInsufficientPermissions.Error())
		}

		ctx = withClaims(ctx, claims)

		return handler(ctx, req)
	}, nil
}

// Context helpers for JWT claims.

type claimsContextKey struct{}

func ClaimsFromContext(ctx context.Context) (*jwt.Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*jwt.Claims)
	return claims, ok
}

func withClaims(ctx context.Context, claims *jwt.Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

// Token extraction from gRPC metadata.

const authorizationMetadataKey = "authorization"

func accessTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, ErrAuthorizationTokenMissing.Error())
	}

	values := md.Get(authorizationMetadataKey)
	for _, value := range values {
		token := extractBearerToken(value)
		if token != "" {
			return token, nil
		}
	}

	return "", status.Error(codes.Unauthenticated, ErrAuthorizationTokenMissing.Error())
}

func extractBearerToken(value string) string {
	value = strings.TrimSpace(value)

	scheme, token, found := strings.Cut(value, " ")
	if found {
		if strings.EqualFold(scheme, "Bearer") {
			return strings.TrimSpace(token)
		}
		// Non-Bearer scheme (e.g., "Basic"), reject.
		return ""
	}

	// No space: reject bare "Bearer", accept raw tokens.
	if strings.EqualFold(value, "Bearer") {
		return ""
	}

	return value
}
