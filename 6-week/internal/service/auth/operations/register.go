package operations

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

func Register(
	ctx context.Context,
	userSvc userService.Service,
	jwtManager jwt.Manager,
	name,
	email,
	password,
	passwordConfirm string,
) (string, domainAuth.TokenPair, error) {
	userID, err := userSvc.Create(ctx, &domainUser.User{
		Name:  name,
		Email: email,
		Role:  domainUser.RoleUser,
	}, password, passwordConfirm)
	if err != nil {
		return "", domainAuth.TokenPair{}, err
	}

	tokens, err := generateTokenPair(jwtManager, userID, int32(domainUser.RoleUser))
	if err != nil {
		return "", domainAuth.TokenPair{}, err
	}

	return userID, tokens, nil
}
