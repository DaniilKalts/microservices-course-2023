package repository

import (
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	userRepository "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user"
	"go.uber.org/zap"
)

type UserRepository = userRepository.Repository

type Repositories struct {
	User UserRepository
}

type Deps struct {
	DB     database.Client
	Logger *zap.Logger
}

func NewRepositories(deps Deps) Repositories {
	return Repositories{
		User: userRepository.NewRepository(deps.DB, deps.Logger.Named("user")),
	}
}
