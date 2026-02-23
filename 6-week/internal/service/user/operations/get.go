package operations

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
)

func Get(ctx context.Context, repo repository.UserRepository, id string) (*domainUser.User, error) {
	return repo.Get(ctx, id)
}
