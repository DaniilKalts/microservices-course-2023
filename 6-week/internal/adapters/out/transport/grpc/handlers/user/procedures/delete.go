package procedures

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type DeleteInput struct {
	ID string
}

func Delete(ctx context.Context, userService service.UserService, input DeleteInput) error {
	return userService.Delete(ctx, input.ID)
}
