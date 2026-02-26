package auth

import (
	"errors"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func mapAuthError(err error) error {
	switch {
	case errors.Is(err, domainAuth.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, domainAuth.ErrInvalidCredentials.Error())
	case errors.Is(err, domainAuth.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, domainAuth.ErrInvalidRefreshToken.Error())
	case errors.Is(err, domainAuth.ErrAuthentication):
		return status.Error(codes.Internal, domainAuth.ErrAuthentication.Error())
	case errors.Is(err, domainAuth.ErrUserIDEmpty):
		return status.Error(codes.Internal, domainAuth.ErrIssueTokens.Error())
	case errors.Is(err, domainAuth.ErrIssueTokens):
		return status.Error(codes.Internal, domainAuth.ErrIssueTokens.Error())
	default:
		return err
	}
}
