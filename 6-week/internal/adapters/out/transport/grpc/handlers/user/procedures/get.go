package procedures

import (
	"context"

	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
	userOperations "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user/operations"
)

func Get(ctx context.Context, userSvc service.UserService, input userOperations.GetInput) (*domainUser.User, error) {
	return userSvc.Get(ctx, input)
}
