package service

import (
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	authService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
	"go.uber.org/zap"
)

type UserService = userService.Service
type AuthService = authService.Service

type Services struct {
	User UserService
	Auth AuthService
}

type Deps struct {
	Repositories repository.Repositories
	JWTManager   jwt.Manager
	Logger       *zap.Logger
}

func NewServices(deps Deps) Services {
	userSvc := userService.NewService(deps.Repositories.User, deps.Logger.Named("user"))
	authSvc := authService.NewService(
		userSvc,
		deps.Repositories.User,
		deps.JWTManager,
		deps.Logger.Named("auth"),
	)

	return Services{
		User: userSvc,
		Auth: authSvc,
	}
}
