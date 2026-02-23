package procedures

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type GetInput struct {
	ID string
}

func Get(ctx context.Context, userService service.UserService, input GetInput) (*domainUser.User, error) {
	return userService.Get(ctx, input.ID)
}
