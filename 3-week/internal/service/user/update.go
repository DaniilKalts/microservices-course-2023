package user

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/models"
)

func (s *service) Update(ctx context.Context, id string, userPatch *models.UpdateUserPatch) error {
	return s.userRepo.Update(ctx, id, userPatch)
}
