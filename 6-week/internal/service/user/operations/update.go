package operations

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
)

func Update(ctx context.Context, repo repository.UserRepository, id string, patch *domainUser.UpdatePatch) error {
	return repo.Update(ctx, id, patch)
}
