package auth

import domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"

type Credentials struct {
	ID           string
	PasswordHash string
	Role         domainUser.Role
}
