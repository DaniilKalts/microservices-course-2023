package user

import (
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

type CreateInput struct {
	User     *domainUser.User
	Password string
}

type UpdateInput struct {
	ID       string
	Name     *string
	Email    *string
	Password *string
}
