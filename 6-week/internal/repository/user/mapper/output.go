package mapper

import (
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user/model"
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
