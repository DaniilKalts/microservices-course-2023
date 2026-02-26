package auth

import "errors"

var (
	ErrUserIDEmpty         = errors.New("user id is empty")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrAuthentication      = errors.New("failed to authenticate user")
	ErrIssueTokens         = errors.New("failed to issue auth tokens")
)
