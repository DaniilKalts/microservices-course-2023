package user

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

func (s *service) Get(ctx context.Context, id string) (*domainUser.User, error) {
	return s.repo.Get(ctx, id)
}
