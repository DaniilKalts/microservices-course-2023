package operations

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
)

type GetInput struct {
	ID string
}

func Get(ctx context.Context, repo repository.UserRepository, input GetInput) (*domainUser.User, error) {
	return repo.Get(ctx, input.ID)
}
