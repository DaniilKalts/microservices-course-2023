package service

import (
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/repository"
	authService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
)

type Services struct {
	User userService.Service
	Auth authService.Service
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
		deps.JWTManager,
		deps.Logger.Named("auth"),
	)

	return Services{
		User: userSvc,
		Auth: authSvc,
	}
}
