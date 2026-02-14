package user

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
)

func (s *service) Get(ctx context.Context, id string) (*domainUser.Entity, error) {
	return s.repo.Get(ctx, id)
}
