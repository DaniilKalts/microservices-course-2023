package auth

import (
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	authService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/auth"
)

func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, authService.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, authService.ErrInvalidCredentials.Error())
	case errors.Is(err, authService.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, authService.ErrInvalidRefreshToken.Error())
	case errors.Is(err, authService.ErrAuthentication):
		return status.Error(codes.Internal, authService.ErrAuthentication.Error())
	case errors.Is(err, authService.ErrUserIDEmpty):
		return status.Error(codes.Internal, authService.ErrUserIDEmpty.Error())
	case errors.Is(err, authService.ErrIssueTokens):
		return status.Error(codes.Internal, authService.ErrIssueTokens.Error())
	case errors.Is(err, domainUser.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, domainUser.ErrEmailAlreadyExists.Error())
	default:
		h.logger.Error("unhandled auth error", zap.Error(err))
		return status.Error(codes.Internal, "internal error")
	}
}
