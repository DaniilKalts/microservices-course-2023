package procedures

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
	authOperations "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/auth/operations"
)

func Login(ctx context.Context, authSvc service.AuthService, input authOperations.LoginInput) (domainAuth.TokenPair, error) {
	return authSvc.Login(ctx, input)
}
