package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

type LogoutInput struct {
	RefreshToken string
}

func Logout(_ context.Context, jwtManager jwt.Manager, input LogoutInput) error {
	return parseRefreshToken(jwtManager, input.RefreshToken)
}
