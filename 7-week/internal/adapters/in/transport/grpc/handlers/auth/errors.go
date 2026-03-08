package auth

import (
	"errors"

	authService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapAuthError(err error) error {
	switch {
	case errors.Is(err, authService.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, authService.ErrInvalidCredentials.Error())
	case errors.Is(err, authService.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, authService.ErrInvalidRefreshToken.Error())
	case errors.Is(err, authService.ErrAuthentication):
		return status.Error(codes.Internal, authService.ErrAuthentication.Error())
	case errors.Is(err, authService.ErrUserIDEmpty):
		return status.Error(codes.Internal, authService.ErrIssueTokens.Error())
	case errors.Is(err, authService.ErrIssueTokens):
		return status.Error(codes.Internal, authService.ErrIssueTokens.Error())
	case errors.Is(err, userService.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, userService.ErrEmailAlreadyExists.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
