package auth

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

func authorize(token string, jwtManager jwt.Manager, requiredGroup AccessGroup) (*jwt.Claims, error) {
	if jwtManager == nil {
		return nil, status.Error(codes.Internal, domainAuth.ErrJWTManagerNotConfigured.Error())
	}

	if strings.TrimSpace(token) == "" {
		return nil, status.Error(codes.Unauthenticated, domainAuth.ErrAuthorizationTokenMissing.Error())
	}

	claims, err := jwtManager.VerifyAccessToken(token)
	if err != nil || claims == nil {
		return nil, status.Error(codes.Unauthenticated, domainAuth.ErrInvalidAccessToken.Error())
	}

	if !hasRequiredRole(claims.RoleID, requiredGroup) {
		return nil, status.Error(codes.PermissionDenied, domainAuth.ErrInsufficientPermissions.Error())
	}

	return claims, nil
}

func hasRequiredRole(roleID int32, requiredGroup AccessGroup) bool {
	switch requiredGroup {
	case AccessGroupAuthenticated:
		return roleID == int32(domainUser.RoleUser) || roleID == int32(domainUser.RoleAdmin)
	case AccessGroupAdmin:
		return roleID == int32(domainUser.RoleAdmin)
	default:
		return false
	}
}
