package user

import "errors"

var (
	ErrNotFound           = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNoFieldsToUpdate   = errors.New("no fields to update")
	ErrWeakPassword       = errors.New("password must contain at least one uppercase letter, one lowercase letter, and one digit")
)
