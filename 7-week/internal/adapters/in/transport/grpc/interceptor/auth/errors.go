package auth

import "errors"

var (
	ErrJWTManagerNotConfigured   = errors.New("jwt manager is not configured")
	ErrAuthorizationTokenMissing = errors.New("authorization token is required")
	ErrInvalidAccessToken        = errors.New("invalid access token")
	ErrInsufficientPermissions   = errors.New("insufficient role permissions")
	ErrAccessPolicyNotConfigured = errors.New("auth access policy is not configured")
	ErrAccessDenied              = errors.New("access denied")
)
