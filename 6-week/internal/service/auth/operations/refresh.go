package operations

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

func Refresh(_ context.Context, jwtManager jwt.Manager, refreshToken string) (domainAuth.TokenPair, error) {
	claims, err := verifyRefreshToken(jwtManager, refreshToken)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	tokens, err := generateTokenPair(jwtManager, claims.ID, claims.RoleID)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	return tokens, nil
}
