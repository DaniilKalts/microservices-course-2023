package user

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
)

func (s *service) Get(ctx context.Context, id string) (*models.User, error) {
	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}