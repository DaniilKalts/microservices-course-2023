package mapper

import (
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/model"
)

func ToDBUserFromDomain(user *domainUser.User, passwordHash string) model.DBUser {
	return model.DBUser{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: passwordHash,
		Role:         int32(user.Role),
	}
}
