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
	// User Service
	userSvc := userService.NewService(deps.Repositories.User)
	userSvc = userService.WithTracing(userSvc)
	userSvc = userService.WithLogging(userSvc, deps.Logger.Named("user"))

	// Auth Service
	authSvc := authService.NewService(userSvc, deps.JWTManager)
	authSvc = authService.WithTracing(authSvc)
	authSvc = authService.WithLogging(authSvc, deps.Logger.Named("auth"))

	return Services{
		User: userSvc,
		Auth: authSvc,
	}
}
