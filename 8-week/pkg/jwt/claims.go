package jwt

import (
	"errors"
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID    string
	RoleID    int32
	TokenType string

	ID        string
	Issuer    string
	Subject   string
	Audience  []string
	ExpiresAt time.Time
	IssuedAt  time.Time
	NotBefore time.Time
}

type jwtClaims struct {
	UserID    string `json:"user_id"`
	RoleID    int32  `json:"role_id"`
	TokenType string `json:"token_type"`
	jwtv5.RegisteredClaims
}

func (c jwtClaims) Validate() error {
	if c.UserID == "" {
		return ErrUserIDMissing
	}

	return validateTokenType(c.TokenType, "")
}

func toJWTClaims(c Claims) jwtClaims {
	jc := jwtClaims{
		UserID:    c.UserID,
		RoleID:    c.RoleID,
		TokenType: c.TokenType,
		RegisteredClaims: jwtv5.RegisteredClaims{
			ID:      c.ID,
			Issuer:  c.Issuer,
			Subject: c.Subject,
		},
	}

	if len(c.Audience) > 0 {
		jc.Audience = jwtv5.ClaimStrings(c.Audience)
	}

	if !c.ExpiresAt.IsZero() {
		jc.ExpiresAt = jwtv5.NewNumericDate(c.ExpiresAt)
	}

	if !c.IssuedAt.IsZero() {
		jc.IssuedAt = jwtv5.NewNumericDate(c.IssuedAt)
	}

	if !c.NotBefore.IsZero() {
		jc.NotBefore = jwtv5.NewNumericDate(c.NotBefore)
	}

	return jc
}

func fromJWTClaims(jc *jwtClaims) *Claims {
	c := &Claims{
		UserID:    jc.UserID,
		RoleID:    jc.RoleID,
		TokenType: jc.TokenType,
		ID:        jc.ID,
		Issuer:    jc.Issuer,
		Subject:   jc.Subject,
		Audience:  []string(jc.Audience),
	}

	if jc.ExpiresAt != nil {
		c.ExpiresAt = jc.ExpiresAt.Time
	}

	if jc.IssuedAt != nil {
		c.IssuedAt = jc.IssuedAt.Time
	}

	if jc.NotBefore != nil {
		c.NotBefore = jc.NotBefore.Time
	}

	return c
}

func (m *manager) prepareClaims(claims Claims, ttl time.Duration, expectedTokenType string) (jwtClaims, error) {
	if err := validateTokenType(claims.TokenType, expectedTokenType); err != nil {
		if !errors.Is(err, ErrTokenTypeMissing) {
			return jwtClaims{}, fmt.Errorf("prepare claims: %w", err)
		}

		claims.TokenType = expectedTokenType
	}

	if claims.UserID == "" {
		return jwtClaims{}, fmt.Errorf("prepare claims: %w", ErrUserIDMissing)
	}

	now := time.Now().UTC()

	if claims.Issuer == "" {
		claims.Issuer = m.issuer
	}

	if len(claims.Audience) == 0 && m.audience != "" {
		claims.Audience = []string{m.audience}
	}

	if claims.Subject == "" {
		claims.Subject = m.subject
	}

	if claims.IssuedAt.IsZero() {
		claims.IssuedAt = now.Add(m.issuedAtOffset)
	}

	if claims.NotBefore.IsZero() {
		claims.NotBefore = now.Add(m.notBeforeOffset)
	}

	if claims.ExpiresAt.IsZero() {
		claims.ExpiresAt = now.Add(ttl)
	}

	claims.ID = uuid.NewString()

	return toJWTClaims(claims), nil
}
