package user

import (
	userv1 "github.com/DaniilKalts/microservices-course-2023/4-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/service"
)

type Handler struct {
	userv1.UnimplementedUserV1Server
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}
