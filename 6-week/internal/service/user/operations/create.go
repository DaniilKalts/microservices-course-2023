package operations

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
)

func Create(
	ctx context.Context,
	repo repository.UserRepository,
	user *domainUser.User,
	password,
	passwordConfirm string,
) (string, error) {
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

	return repo.Create(ctx, user, string(hashedPassword))
}
