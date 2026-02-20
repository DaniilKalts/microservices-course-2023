package auth

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
)

func (s *service) Refresh(_ context.Context, refreshToken string) (domainAuth.TokenPair, error) {
	claims, err := s.verifyRefreshToken(refreshToken)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	tokens, err := s.generateTokenPair(claims.ID, claims.RoleID)
	if err != nil {
		return domainAuth.TokenPair{}, err
	}

	return tokens, nil
}
