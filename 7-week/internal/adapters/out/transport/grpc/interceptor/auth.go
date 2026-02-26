package interceptor

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

const (
	authorizationMetadataKey = "authorization"
)

var methodAccessPolicy = map[string]domainUser.Role{
	userv1.UserV1_Create_FullMethodName: domainUser.RoleAdmin,
	userv1.UserV1_Update_FullMethodName: domainUser.RoleAdmin,
	userv1.UserV1_Delete_FullMethodName: domainUser.RoleAdmin,
}

func AuthInterceptor(jwtManager jwt.Manager) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requiredRole, requiresAuth := requiredRole(info.FullMethod)
		if !requiresAuth {
			return handler(ctx, req)
		}

		token, err := accessTokenFromContext(ctx)
		if err != nil {
			return nil, err
		}

		if err = authorize(token, jwtManager, requiredRole); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func requiredRole(fullMethod string) (domainUser.Role, bool) {
	role, ok := methodAccessPolicy[fullMethod]
	return role, ok
}

func accessTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "authorization token is required")
	}

	values := md.Get(authorizationMetadataKey)
	for _, value := range values {
		if token := strings.TrimSpace(value); token != "" {
			return token, nil
		}
	}

	return "", status.Error(codes.Unauthenticated, "authorization token is required")
}

func authorize(token string, jwtManager jwt.Manager, requiredRole domainUser.Role) error {
	if jwtManager == nil {
		return status.Error(codes.Internal, "jwt manager is not configured")
	}

	if strings.TrimSpace(token) == "" {
		return status.Error(codes.Unauthenticated, "authorization token is required")
	}

	claims, err := jwtManager.VerifyAccessToken(token)
	if err != nil || claims == nil {
		return status.Error(codes.Unauthenticated, "invalid access token")
	}

	if !hasRequiredRole(claims.RoleID, requiredRole) {
		return status.Error(codes.PermissionDenied, "insufficient role permissions")
	}

	return nil
}

func hasRequiredRole(roleID int32, requiredRole domainUser.Role) bool {
	switch requiredRole {
	case domainUser.RoleUser:
		return roleID == int32(domainUser.RoleUser) || roleID == int32(domainUser.RoleAdmin)
	case domainUser.RoleAdmin:
		return roleID == int32(domainUser.RoleAdmin)
	default:
		return false
	}
}
