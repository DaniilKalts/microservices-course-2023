package auth

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/auth/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	authService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/auth"
)

// Proto → Domain

func toRegisterInput(req *authv1.RegisterRequest) authService.RegisterInput {
	return authService.RegisterInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func toLoginInput(req *authv1.LoginRequest) authService.LoginInput {
	return authService.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func toLogoutInput(req *authv1.LogoutRequest) authService.LogoutInput {
	return authService.LogoutInput{RefreshToken: req.GetRefreshToken()}
}

func toRefreshInput(req *authv1.RefreshRequest) authService.RefreshInput {
	return authService.RefreshInput{RefreshToken: req.GetRefreshToken()}
}

// Domain → Proto

func toProtoTokenPair(tokens authService.TokenPair) *authv1.TokenPair {
	return &authv1.TokenPair{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresIn:  tokens.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: tokens.RefreshTokenExpiresIn,
	}
}

func toProtoRegisterUser(user domainUser.User) *authv1.RegisterResponse_User {
	return &authv1.RegisterResponse_User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
