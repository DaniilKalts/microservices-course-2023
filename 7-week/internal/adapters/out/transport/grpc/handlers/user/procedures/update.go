package procedures

import (
	"context"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	userOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user/operations"
)

func Update(ctx context.Context, userSvc service.UserService, input userOperations.UpdateInput) error {
	return userSvc.Update(ctx, input)
}
