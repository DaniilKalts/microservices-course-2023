package service

import (
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	authService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
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
}

func NewServices(deps Deps) Services {
	userSvc := userService.NewService(deps.Repositories.User)
	authSvc := authService.NewService(userSvc, deps.Repositories.User, deps.JWTManager)

	return Services{
		User: userSvc,
		Auth: authSvc,
	}
}
