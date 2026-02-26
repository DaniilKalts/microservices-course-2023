package repository

import (
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	userRepository "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user"
)

type UserRepository = userRepository.Repository

type Repositories struct {
	User UserRepository
}

type Deps struct {
	DB database.Client
}

func NewRepositories(deps Deps) Repositories {
	return Repositories{
		User: userRepository.NewRepository(deps.DB),
	}
}
