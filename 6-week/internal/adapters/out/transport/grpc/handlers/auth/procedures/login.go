package procedures

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type LoginInput struct {
	Email    string
	Password string
}

func Login(ctx context.Context, authService service.AuthService, input LoginInput) (domainAuth.TokenPair, error) {
	return authService.Login(ctx, input.Email, input.Password)
}
