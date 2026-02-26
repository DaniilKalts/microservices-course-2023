package operations

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
)

type CreateInput struct {
	User            *domainUser.User
	Password        string
	PasswordConfirm string
}

func Create(
	ctx context.Context,
	repo repository.UserRepository,
	input CreateInput,
) (string, error) {
	if input.Password != input.PasswordConfirm {
		return "", errors.New("Passwords don't match")
	}

	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	input.User.ID = id.String()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return repo.Create(ctx, input.User, string(hashedPassword))
}
