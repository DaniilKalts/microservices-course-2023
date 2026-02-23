package procedures

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type RefreshInput struct {
	RefreshToken string
}

func Refresh(ctx context.Context, authService service.AuthService, input RefreshInput) (domainAuth.TokenPair, error) {
	return authService.Refresh(ctx, input.RefreshToken)
}
