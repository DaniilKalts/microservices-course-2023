package auth

import "errors"

var (
	ErrUserIDEmpty        = errors.New("user id is empty")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
