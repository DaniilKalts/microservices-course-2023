package operations

import "errors"

var (
	errInvalidTokenType    = errors.New("invalid token type")
	errRefreshTokenEmpty   = errors.New("refresh token is empty")
	errUserIDEmpty         = errors.New("user id is empty")
	errLoginNotImplemented = errors.New("login requires credentials verification flow")
)
