package auth

import (
	"strings"

	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

const bearerScheme = "Bearer"

func normalizeRefreshToken(refreshToken string) string {
	refreshToken = strings.TrimSpace(refreshToken)
	if refreshToken == "" {
		return ""
	}

	parts := strings.Fields(refreshToken)
	if len(parts) == 2 && strings.EqualFold(parts[0], bearerScheme) {
		return parts[1]
	}

	if len(parts) == 1 && strings.EqualFold(parts[0], bearerScheme) {
		return ""
	}

	return refreshToken
}

func (s *service) verifyRefreshToken(refreshToken string) (*jwt.Claims, error) {
	normalized := normalizeRefreshToken(refreshToken)
	if normalized == "" {
		return nil, errRefreshTokenEmpty
	}

	claims, err := s.jwtManager.Verify(normalized)
	if err != nil {
		return nil, err
	}

	if claims.TokenType != refreshTokenType {
		return nil, errInvalidTokenType
	}

	if claims.ID == "" {
		return nil, errUserIDEmpty
	}

	return claims, nil
}

func (s *service) parseRefreshToken(refreshToken string) error {
	normalized := normalizeRefreshToken(refreshToken)
	if normalized == "" {
		return errRefreshTokenEmpty
	}

	claims, err := s.jwtManager.Parse(normalized)
	if err != nil {
		return err
	}

	if claims.TokenType != refreshTokenType {
		return errInvalidTokenType
	}

	if claims.ID == "" {
		return errUserIDEmpty
	}

	return nil
}
