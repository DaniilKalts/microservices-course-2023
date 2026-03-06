package auth

import "errors"

var (
	ErrUserIDEmpty               = errors.New("user id is empty")
	ErrInvalidCredentials        = errors.New("invalid email or password")
	ErrInvalidRefreshToken       = errors.New("invalid refresh token")
	ErrAuthentication            = errors.New("failed to authenticate user")
	ErrIssueTokens               = errors.New("failed to issue auth tokens")
	ErrJWTManagerNotConfigured   = errors.New("jwt manager is not configured")
	ErrAuthorizationTokenMissing = errors.New("authorization token is required")
	ErrInvalidAccessToken        = errors.New("invalid access token")
	ErrInsufficientPermissions   = errors.New("insufficient role permissions")
	ErrAccessPolicyNotConfigured = errors.New("auth access policy is not configured")
	ErrAccessDenied              = errors.New("access denied")
)
