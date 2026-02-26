package jwt

import (
	"errors"
	"fmt"
	"strings"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
	bearerScheme     = "Bearer"
)

var (
	errTokenEmpty        = errors.New("token is empty")
	errTokenIDMissing    = errors.New("token id is missing")
	errTokenTypeMissing  = errors.New("token type is missing")
	errTokenTypeInvalid  = errors.New("invalid token type")
	errTokenTypeMismatch = errors.New("token type mismatch")
)

func validateTokenType(tokenType, expectedTokenType string) error {
	if tokenType == "" {
		return errTokenTypeMissing
	}

	switch tokenType {
	case tokenTypeAccess, tokenTypeRefresh:
	default:
		return fmt.Errorf("%w: %q", errTokenTypeInvalid, tokenType)
	}

	if expectedTokenType != "" && tokenType != expectedTokenType {
		return fmt.Errorf("%w: got %q, want %q", errTokenTypeMismatch, tokenType, expectedTokenType)
	}

	return nil
}

func normalizeToken(tokenString string) string {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return ""
	}

	parts := strings.Fields(tokenString)
	if len(parts) == 1 && strings.EqualFold(parts[0], bearerScheme) {
		return ""
	}

	if len(parts) == 2 && strings.EqualFold(parts[0], bearerScheme) {
		return parts[1]
	}

	return tokenString
}

func (m *manager) publicKeyFunc(_ *jwtv5.Token) (any, error) {
	return m.publicKey, nil
}
