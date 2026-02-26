package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
)

type DeleteInput struct {
	ID string
}

func Delete(ctx context.Context, repo repository.UserRepository, input DeleteInput) error {
	return repo.Delete(ctx, input.ID)
}
