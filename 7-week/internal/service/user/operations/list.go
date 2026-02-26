package operations

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
)

func List(ctx context.Context, repo repository.UserRepository) ([]domainUser.User, error) {
	return repo.List(ctx)
}
