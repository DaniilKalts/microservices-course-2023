package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	userRepositoryOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository/user/operations"
)

type UpdateInput struct {
	ID    string
	Name  *string
	Email *string
}

func Update(ctx context.Context, repo repository.UserRepository, input UpdateInput) error {
	return repo.Update(ctx, userRepositoryOperations.UpdateInput{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
}
