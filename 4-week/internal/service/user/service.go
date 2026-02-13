package user

import (
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/repository"
	srv "github.com/DaniilKalts/microservices-course-2023/4-week/internal/service"
)

type service struct {
	repo repository.UserRepository
}

func NewService(repo repository.UserRepository) srv.UserService {
	return &service{repo: repo}
}
