package user

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

func (s *service) Update(ctx context.Context, id string, patch *domainUser.UpdatePatch) error {
	return s.repo.Update(ctx, id, patch)
}
