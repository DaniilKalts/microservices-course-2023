package user

import domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"

func toDBUserFromDomain(user *domainUser.User) *dbUser {
	return &dbUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  int32(user.Role),
	}
}

func toDomainFromDBUser(user *dbUser) *domainUser.User {
	return &domainUser.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      domainUser.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
