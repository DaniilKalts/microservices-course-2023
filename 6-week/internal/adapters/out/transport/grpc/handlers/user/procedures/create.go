package procedures

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type CreateInput struct {
	User            *domainUser.User
	Password        string
	PasswordConfirm string
}

func Create(ctx context.Context, userService service.UserService, input CreateInput) (string, error) {
	return userService.Create(ctx, input.User, input.Password, input.PasswordConfirm)
}
