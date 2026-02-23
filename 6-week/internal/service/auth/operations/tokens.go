package operations

import (
	jwtv5 "github.com/golang-jwt/jwt/v5"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

func generateTokenPair(jwtManager jwt.Manager, userID string, roleID int32) (domainAuth.TokenPair, error) {
	if userID == "" {
		return domainAuth.TokenPair{}, errUserIDEmpty
	}

	claims := jwt.Claims{
		RoleID: roleID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID: userID,
		},
	}

	accessToken, err := jwtManager.GenerateAccessToken(claims)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	refreshToken, err := jwtManager.GenerateRefreshToken(claims)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	return domainAuth.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
