package procedures

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	authOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth/operations"
)

func Refresh(ctx context.Context, authSvc service.AuthService, input authOperations.RefreshInput) (domainAuth.TokenPair, error) {
	return authSvc.Refresh(ctx, input)
}
