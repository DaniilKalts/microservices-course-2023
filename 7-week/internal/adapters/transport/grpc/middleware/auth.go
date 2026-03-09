package middleware

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

var (
	ErrJWTManagerNotConfigured   = errors.New("jwt manager is not configured")
	ErrAuthorizationTokenMissing = errors.New("authorization token is required")
	ErrInvalidAccessToken        = errors.New("invalid access token")
	ErrInsufficientPermissions   = errors.New("insufficient role permissions")
	ErrAccessPolicyNotConfigured = errors.New("auth access policy is not configured")
	ErrMethodNotRegistered       = errors.New("method not registered in access policy")
)

func AuthInterceptor(jwtManager jwt.Manager, policy AccessPolicy) grpc.UnaryServerInterceptor {
	if policy.IsEmpty() {
		return func(
			_ context.Context,
			_ any,
			_ *grpc.UnaryServerInfo,
			_ grpc.UnaryHandler,
		) (any, error) {
			return nil, status.Error(codes.Internal, ErrAccessPolicyNotConfigured.Error())
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
			return nil, status.Error(codes.PermissionDenied, ErrMethodNotRegistered.Error())
		}

		if requiredGroup.IsPublic {
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

func authorize(token string, jwtManager jwt.Manager, requiredGroup AccessGroup) (*jwt.Claims, error) {
	if jwtManager == nil {
		return nil, status.Error(codes.Internal, ErrJWTManagerNotConfigured.Error())
	}

	if strings.TrimSpace(token) == "" {
		return nil, status.Error(codes.Unauthenticated, ErrAuthorizationTokenMissing.Error())
	}

	claims, err := jwtManager.VerifyAccessToken(token)
	if err != nil || claims == nil {
		return nil, status.Error(codes.Unauthenticated, ErrInvalidAccessToken.Error())
	}

	if !requiredGroup.AllowsRole(claims.RoleID) {
		return nil, status.Error(codes.PermissionDenied, ErrInsufficientPermissions.Error())
	}

	return claims, nil
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
