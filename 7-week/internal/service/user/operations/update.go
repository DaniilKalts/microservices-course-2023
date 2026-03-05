package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	userRepositoryOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/operations"
	"golang.org/x/crypto/bcrypt"
)

type UpdateInput struct {
	ID       string
	Name     *string
	Email    *string
	Password *string
}

func Update(ctx context.Context, repo repository.UserRepository, input UpdateInput) error {
	repoInput := userRepositoryOperations.UpdateInput{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	}

	if input.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		passwordHash := string(hashedPassword)
		repoInput.PasswordHash = &passwordHash
	}

	return repo.Update(ctx, repoInput)
}
