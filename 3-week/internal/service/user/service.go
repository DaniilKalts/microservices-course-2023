package user

import (
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository"
	srv "github.com/DaniilKalts/microservices-course-2023/3-week/internal/service"
)

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) srv.UserService {
	return &service{
		userRepo: userRepo,
	}
}