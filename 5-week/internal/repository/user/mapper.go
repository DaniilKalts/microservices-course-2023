package user

import domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"

func toDBUserFromDomain(user *domainUser.Entity) *dbUser {
	return &dbUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  int32(user.Role),
	}
}

func toDomainFromDBUser(user *dbUser) *domainUser.Entity {
	updatedAt := user.UpdatedAt

	return &domainUser.Entity{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      domainUser.Role(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: &updatedAt,
	}
}
