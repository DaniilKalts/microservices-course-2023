package operations

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

type LoginInput struct {
	Email    string
	Password string
}

const dummyPasswordHash = "$2a$10$7EqJtq98hPqEX7fNZaFWoO5m6jH9MuNoMNFcQJUO2cMjJ1ytY1/6W"

func Login(
	ctx context.Context,
	userRepo repository.UserRepository,
	jwtManager jwt.Manager,
	input LoginInput,
) (domainAuth.TokenPair, error) {
	credentials, err := userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(input.Password))
		if pgxscan.NotFound(err) {
			return domainAuth.TokenPair{}, status.Error(codes.Unauthenticated, errInvalidCredentials.Error())
		}

		return domainAuth.TokenPair{}, status.Error(codes.Internal, "failed to authenticate user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentials.PasswordHash), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domainAuth.TokenPair{}, status.Error(codes.Unauthenticated, errInvalidCredentials.Error())
		}

		return domainAuth.TokenPair{}, status.Error(codes.Internal, "failed to authenticate user")
	}

	tokens, err := generateTokenPair(jwtManager, credentials.ID, int32(credentials.Role))
	if err != nil {
		return domainAuth.TokenPair{}, status.Error(codes.Internal, "failed to issue auth tokens")
	}

	return tokens, nil
}
