package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type LogoutInput struct {
	RefreshToken string
}

func Logout(_ context.Context, jwtManager jwt.Manager, input LogoutInput) error {
	_, err := jwtManager.VerifyRefreshToken(input.RefreshToken)
	return err
}
