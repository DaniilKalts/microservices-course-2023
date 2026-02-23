package procedures

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type UpdateInput struct {
	ID    string
	Patch *domainUser.UpdatePatch
}

func Update(ctx context.Context, userService service.UserService, input UpdateInput) error {
	return userService.Update(ctx, input.ID, input.Patch)
}
