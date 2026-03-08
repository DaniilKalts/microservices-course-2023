package profile

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
)

func mapProfileError(err error) error {
	switch {
	case errors.Is(err, userService.ErrNotFound):
		return status.Error(codes.NotFound, userService.ErrNotFound.Error())
	case errors.Is(err, userService.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, userService.ErrEmailAlreadyExists.Error())
	case errors.Is(err, userService.ErrNoFieldsToUpdate):
		return status.Error(codes.InvalidArgument, userService.ErrNoFieldsToUpdate.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
