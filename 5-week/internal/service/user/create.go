package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
)

func (s *service) Create(ctx context.Context, user *domainUser.Entity, password, passwordConfirm string) (string, error) {
	if password != passwordConfirm {
		return "", errors.New("Passwords don't match")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	user.ID = id.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return s.repo.Create(ctx, user, string(hashedPassword))
}
