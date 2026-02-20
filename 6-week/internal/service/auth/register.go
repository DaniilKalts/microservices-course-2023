package auth

import (
	"context"

	domainAuth "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

func (s *service) Register(ctx context.Context, name, email, password, passwordConfirm string) (string, domainAuth.TokenPair, error) {
	userID, err := s.userService.Create(ctx, &domainUser.User{
		Name:  name,
		Email: email,
		Role:  domainUser.RoleUser,
	}, password, passwordConfirm)
	if err != nil {
		return "", domainAuth.TokenPair{}, err
	}

	tokens, err := s.generateTokenPair(userID, int32(domainUser.RoleUser))
	if err != nil {
		return "", domainAuth.TokenPair{}, err
	}

	return userID, tokens, nil
}
