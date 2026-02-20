package auth

import (
	srv "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

const refreshTokenType = "refresh"

type service struct {
	userService srv.UserService
	jwtManager  jwt.Manager
}

func NewService(userService srv.UserService, jwtManager jwt.Manager) srv.AuthService {
	return &service{
		userService: userService,
		jwtManager:  jwtManager,
	}
}
