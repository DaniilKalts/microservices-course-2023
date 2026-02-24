package operations

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

type RefreshInput struct {
	RefreshToken string
}

func Refresh(_ context.Context, jwtManager jwt.Manager, input RefreshInput) (domainAuth.TokenPair, error) {
	claims, err := verifyRefreshToken(jwtManager, input.RefreshToken)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	tokens, err := generateTokenPair(jwtManager, claims.ID, claims.RoleID)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	return tokens, nil
}
