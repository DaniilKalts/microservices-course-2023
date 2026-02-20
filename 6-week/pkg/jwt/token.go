package jwt

import (
	"errors"
	"fmt"
	"strings"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenType  = "access"
	refreshTokenType = "refresh"
	bearerScheme     = "Bearer"
)

var errTokenTypeMissing = errors.New("token type is missing")

func validateTokenType(tokenType, expectedTokenType string) error {
	if tokenType == "" {
		return errTokenTypeMissing
	}

	switch tokenType {
	case accessTokenType, refreshTokenType:
	default:
		return fmt.Errorf("invalid token type %q", tokenType)
	}

	if expectedTokenType != "" && tokenType != expectedTokenType {
		return fmt.Errorf("token type mismatch: got %q", tokenType)
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

func (m *manager) keyFunc(token *jwtv5.Token) (any, error) {
	if token == nil || token.Method == nil {
		return nil, errors.New("token method is missing")
	}

	if token.Method.Alg() != signingMethodAlgorithm {
		return nil, fmt.Errorf("unexpected signing method %q", token.Method.Alg())
	}

	return m.publicKey, nil
}
