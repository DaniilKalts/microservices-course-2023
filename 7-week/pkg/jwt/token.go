package jwt

import (
	"fmt"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const (
	tokenTypeAccess  = "access"
	tokenTypeRefresh = "refresh"
)

func validateTokenType(tokenType, expectedTokenType string) error {
	if tokenType == "" {
		return ErrTokenTypeMissing
	}

	switch tokenType {
	case tokenTypeAccess, tokenTypeRefresh:
	default:
		return fmt.Errorf("%w: %q", ErrTokenTypeInvalid, tokenType)
	}

	if expectedTokenType != "" && tokenType != expectedTokenType {
		return fmt.Errorf("%w: got %q, want %q", ErrTokenTypeMismatch, tokenType, expectedTokenType)
	}

	return nil
}

func (m *manager) publicKeyFunc(token *jwtv5.Token) (any, error) {
	if _, ok := token.Method.(*jwtv5.SigningMethodRSA); !ok {
		return nil, fmt.Errorf("%w: %v", ErrSigningMethod, token.Header["alg"])
	}

	return m.publicKey, nil
}
