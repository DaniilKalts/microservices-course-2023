package operations

import "errors"

var (
	errInvalidTokenType   = errors.New("invalid token type")
	errRefreshTokenEmpty  = errors.New("refresh token is empty")
	errUserIDEmpty        = errors.New("user id is empty")
	errInvalidCredentials = errors.New("invalid email or password")
)
