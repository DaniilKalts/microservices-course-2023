package operations

import (
	"context"
	"fmt"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
	userOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user/operations"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func Register(
	ctx context.Context,
	userSvc userService.Service,
	jwtManager jwt.Manager,
	input RegisterInput,
) (domainUser.User, domainAuth.TokenPair, error) {
	userID, err := userSvc.Create(ctx, userOperations.CreateInput{
		User: &domainUser.User{
			Name:  input.Name,
			Email: input.Email,
			Role:  domainUser.RoleUser,
		},
		Password:        input.Password,
		PasswordConfirm: input.Password,
	})
	if err != nil {
		return domainUser.User{}, domainAuth.TokenPair{}, err
	}

	tokenPair, err := generateTokenPair(jwtManager, userID, int32(domainUser.RoleUser))
	if err != nil {
		return domainUser.User{}, domainAuth.TokenPair{}, fmt.Errorf("%w: %v", domainAuth.ErrIssueTokens, err)
	}

	return domainUser.User{ID: userID, Name: input.Name, Email: input.Email}, tokenPair, nil
}
