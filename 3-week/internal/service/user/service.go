package user

import (
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository"
	srv "github.com/DaniilKalts/microservices-course-2023/3-week/internal/service"
)

type service struct {
	userRepo  repository.UserRepository
	txManager database.TxManager
}

func NewService(userRepo repository.UserRepository, txManager database.TxManager) srv.UserService {
	return &service{
		userRepo:  userRepo,
		txManager: txManager,
	}
}
