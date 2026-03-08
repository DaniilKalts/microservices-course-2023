package user

import (
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
)

func toDBUser(user *domainUser.User, passwordHash string) dbUser {
	return dbUser{
		ID:           user.ID,
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: passwordHash,
		Role:         int32(user.Role),
	}
}

func toDomainUser(u *dbUser) *domainUser.User {
	return &domainUser.User{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Role:      domainUser.Role(u.Role),
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func toDomainUsers(users []dbUser) []domainUser.User {
	result := make([]domainUser.User, len(users))
	for i := range users {
		result[i] = *toDomainUser(&users[i])
	}
	return result
}

func toCredentials(u *dbUser) *domainUser.Credentials {
	return &domainUser.Credentials{
		ID:           u.ID,
		PasswordHash: u.PasswordHash,
		Role:         domainUser.Role(u.Role),
	}
}
