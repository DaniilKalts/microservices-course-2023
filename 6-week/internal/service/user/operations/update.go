package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
)

type UpdateInput struct {
	ID    string
	Name  *string
	Email *string
}

func Update(ctx context.Context, repo repository.UserRepository, input UpdateInput) error {
	return repo.Update(ctx, repository.UserUpdateInput{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
}
