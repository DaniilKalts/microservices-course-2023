package user

import (
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/repository"
)

type Repository struct {
	dbc database.Client
}

func NewRepository(dbc database.Client) repository.UserRepository {
	return &Repository{dbc: dbc}
}
