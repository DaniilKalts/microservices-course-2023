package procedures

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type RegisterInput struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
}

func Register(
	ctx context.Context,
	authService service.AuthService,
	input RegisterInput,
) (string, domainAuth.TokenPair, error) {
	return authService.Register(ctx, input.Name, input.Email, input.Password, input.PasswordConfirm)
}
