package procedures

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
)

func List(ctx context.Context, userService service.UserService) ([]domainUser.User, error) {
	return userService.List(ctx)
}
