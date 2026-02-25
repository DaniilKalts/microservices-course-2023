package operations

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
	userOperations "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user/operations"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
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

	tokens, err := generateTokenPair(jwtManager, userID, int32(domainUser.RoleUser))
	if err != nil {
		return domainUser.User{}, domainAuth.TokenPair{}, err
	}

	return domainUser.User{ID: userID, Name: input.Name, Email: input.Email}, tokens, nil
}
