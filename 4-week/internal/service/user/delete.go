package user

import (
	"context"
)

func (s *service) Delete(ctx context.Context, id string,) error {
	return s.userRepo.Delete(ctx, id)
}
