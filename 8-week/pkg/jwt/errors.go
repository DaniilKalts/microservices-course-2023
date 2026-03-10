package jwt

import "errors"

var (
	ErrTokenEmpty        = errors.New("token is empty")
	ErrTokenIDMissing    = errors.New("token id is missing")
	ErrTokenTypeMissing  = errors.New("token type is missing")
	ErrTokenTypeInvalid  = errors.New("invalid token type")
	ErrTokenTypeMismatch = errors.New("token type mismatch")
	ErrUserIDMissing     = errors.New("user id is missing")
	ErrSigningMethod     = errors.New("unexpected signing method")
)
