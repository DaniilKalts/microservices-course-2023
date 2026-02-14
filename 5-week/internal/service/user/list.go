package user

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/5-week/internal/domain/user"
)

func (s *service) List(ctx context.Context) ([]domainUser.Entity, error) {
	return s.repo.List(ctx)
}
