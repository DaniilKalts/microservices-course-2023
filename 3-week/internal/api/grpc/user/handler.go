package user

import (
	userv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/service"
)

type Implementation struct {
	userv1.UnimplementedUserV1Server
	userService service.UserService
}

func NewImplementation(userService service.UserService) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
