package profile

import (
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, domainUser.ErrNotFound):
		return status.Error(codes.NotFound, domainUser.ErrNotFound.Error())
	case errors.Is(err, domainUser.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, domainUser.ErrEmailAlreadyExists.Error())
	case errors.Is(err, domainUser.ErrNoFieldsToUpdate):
		return status.Error(codes.InvalidArgument, domainUser.ErrNoFieldsToUpdate.Error())
	default:
		h.logger.Error("unhandled profile error", zap.Error(err))
		return status.Error(codes.Internal, "internal error")
	}
}
