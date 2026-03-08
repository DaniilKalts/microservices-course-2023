package mapper

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	authService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth"
)

func ToRegisterInput(req *authv1.RegisterRequest) authService.RegisterInput {
	return authService.RegisterInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func ToLoginInput(req *authv1.LoginRequest) authService.LoginInput {
	return authService.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func ToLogoutInput(req *authv1.LogoutRequest) authService.LogoutInput {
	return authService.LogoutInput{RefreshToken: req.GetRefreshToken()}
}

func ToRefreshInput(req *authv1.RefreshRequest) authService.RefreshInput {
	return authService.RefreshInput{RefreshToken: req.GetRefreshToken()}
}
