package operations

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
)

func Delete(ctx context.Context, repo repository.UserRepository, id string) error {
	return repo.Delete(ctx, id)
}
