package repository

import (
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	userRepository "github.com/DaniilKalts/microservices-course-2023/8-week/internal/repository/user"
)

type Repositories struct {
	User userRepository.Repository
}

func NewRepositories(db database.Client, logger *zap.Logger) Repositories {
	return Repositories{
		User: userRepository.NewRepository(db, logger.Named("user")),
	}
}
