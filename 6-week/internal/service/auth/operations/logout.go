package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

func Logout(_ context.Context, jwtManager jwt.Manager, refreshToken string) error {
	return parseRefreshToken(jwtManager, refreshToken)
}
