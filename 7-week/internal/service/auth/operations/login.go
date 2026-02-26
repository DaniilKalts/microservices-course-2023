package operations

import (
	"context"
	"errors"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"golang.org/x/crypto/bcrypt"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
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
	credentials, err := userRepo.GetCredentialsByEmail(ctx, input.Email)
	if err != nil {
		_ = bcrypt.CompareHashAndPassword([]byte(dummyPasswordHash), []byte(input.Password))
		if pgxscan.NotFound(err) {
			return domainAuth.TokenPair{}, domainAuth.ErrInvalidCredentials
		}

		return domainAuth.TokenPair{}, fmt.Errorf("%w: %v", domainAuth.ErrAuthentication, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(credentials.PasswordHash), []byte(input.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return domainAuth.TokenPair{}, domainAuth.ErrInvalidCredentials
		}

		return domainAuth.TokenPair{}, fmt.Errorf("%w: %v", domainAuth.ErrAuthentication, err)
	}

	tokenPair, err := generateTokenPair(jwtManager, credentials.ID, int32(credentials.Role))
	if err != nil {
		return domainAuth.TokenPair{}, fmt.Errorf("%w: %v", domainAuth.ErrIssueTokens, err)
	}

	return tokenPair, nil
}
