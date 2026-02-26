package procedures

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	authOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth/operations"
)

func Logout(ctx context.Context, authSvc service.AuthService, input authOperations.LogoutInput) error {
	return authSvc.Logout(ctx, input)
}
