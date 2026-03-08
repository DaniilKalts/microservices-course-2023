package auth

import (
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

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
