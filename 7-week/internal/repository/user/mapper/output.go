package mapper

import (
	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/model"
)

func ToDomainFromDBUser(user *model.DBUser) *domainUser.User {
	return &domainUser.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      domainUser.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func ToDomainFromDBUsers(users []model.DBUser) []domainUser.User {
	entities := make([]domainUser.User, 0, len(users))
	for i := range users {
		entities = append(entities, *ToDomainFromDBUser(&users[i]))
	}

	return entities
}

func ToCredentialsFromDBUser(user *model.DBUser) *domainAuth.Credentials {
	return &domainAuth.Credentials{
		ID:           user.ID,
		PasswordHash: user.PasswordHash,
		Role:         domainUser.Role(user.Role),
	}
}
