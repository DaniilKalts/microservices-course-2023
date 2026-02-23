package procedures

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type LogoutInput struct {
	RefreshToken string
}

func Logout(ctx context.Context, authService service.AuthService, input LogoutInput) error {
	return authService.Logout(ctx, input.RefreshToken)
}
