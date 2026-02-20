package auth

import (
	jwtv5 "github.com/golang-jwt/jwt/v5"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

func (s *service) generateTokenPair(userID string, roleID int32) (domainAuth.TokenPair, error) {
	if userID == "" {
		return domainAuth.TokenPair{}, errUserIDEmpty
	}

	claims := jwt.Claims{
		RoleID: roleID,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID: userID,
		},
	}

	accessToken, err := s.jwtManager.GenerateAccessToken(claims)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(claims)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	return domainAuth.TokenPair{AccessToken: accessToken, RefreshToken: refreshToken}, nil
}
