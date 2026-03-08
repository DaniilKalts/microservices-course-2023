package mapper

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	authService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth"
)

func ToProtoTokenPair(tokens authService.TokenPair) *authv1.TokenPair {
	return &authv1.TokenPair{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresIn:  tokens.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: tokens.RefreshTokenExpiresIn,
	}
}

func ToProtoRegisterUser(user domainUser.User) *authv1.RegisterResponse_User {
	return &authv1.RegisterResponse_User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
