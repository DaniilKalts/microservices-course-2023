package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID    string `json:"user_id"`
	RoleID    int32  `json:"role_id"`
	TokenType string `json:"token_type"`
	jwtv5.RegisteredClaims
}

func (m *manager) prepareClaims(claims Claims, ttl time.Duration, expectedTokenType string) (Claims, error) {
	if err := validateTokenType(claims.TokenType, expectedTokenType); err != nil {
		if !errors.Is(err, errTokenTypeMissing) {
			return Claims{}, fmt.Errorf("prepare claims: %w", err)
		}

		claims.TokenType = expectedTokenType
	}

	now := time.Now().UTC()

	if claims.Issuer == "" {
		claims.Issuer = m.issuer
	}

	if len(claims.Audience) == 0 && m.audience != "" {
		claims.Audience = jwtv5.ClaimStrings{m.audience}
	}

	if claims.Subject == "" {
		claims.Subject = m.subject
	}

	if claims.IssuedAt == nil {
		claims.IssuedAt = jwtv5.NewNumericDate(now.Add(m.issuedAtOffset))
	}

	if claims.NotBefore == nil {
		claims.NotBefore = jwtv5.NewNumericDate(now.Add(m.notBeforeOffset))
	}

	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwtv5.NewNumericDate(now.Add(ttl))
	}

	claims.ID = uuid.NewString()

	return claims, nil
}
