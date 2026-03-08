package service

import (
	authService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
	"go.uber.org/zap"
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
	userSvc := userService.NewService(deps.Repositories.User)
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
