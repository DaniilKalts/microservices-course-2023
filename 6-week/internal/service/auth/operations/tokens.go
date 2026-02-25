package operations

import (
	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

func generateTokenPair(jwtManager jwt.Manager, userID string, roleID int32) (domainAuth.TokenPair, error) {
	if userID == "" {
		return domainAuth.TokenPair{}, domainAuth.ErrUserIDEmpty
	}

	claims := jwt.Claims{
		UserID: userID,
		RoleID: roleID,
	}

	accessToken, err := jwtManager.GenerateAccessToken(claims)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	refreshToken, err := jwtManager.GenerateRefreshToken(claims)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	accessTokenExpiresIn := int64(jwtManager.AccessTokenTTL().Seconds())
	refreshTokenExpiresIn := int64(jwtManager.RefreshTokenTTL().Seconds())

	return domainAuth.TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		AccessTokenExpiresIn:  accessTokenExpiresIn,
		RefreshTokenExpiresIn: refreshTokenExpiresIn,
	}, nil
}
