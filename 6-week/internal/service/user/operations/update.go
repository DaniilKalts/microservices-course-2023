package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	userRepository "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user"
)

type UpdateInput struct {
	ID    string
	Name  *string
	Email *string
}

func Update(ctx context.Context, repo repository.UserRepository, input UpdateInput) error {
	return repo.Update(ctx, userRepository.UpdateInput{
		ID:    input.ID,
		Name:  input.Name,
		Email: input.Email,
	})
}
