package auth

import (
	"context"
)

func (s *service) Logout(_ context.Context, refreshToken string) error {
	return s.parseRefreshToken(refreshToken)
}
