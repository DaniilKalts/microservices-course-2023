package mapper

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
	authOperations "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/auth/operations"
)

func ToRegisterInput(req *authv1.RegisterRequest) authOperations.RegisterInput {
	return authOperations.RegisterInput{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func ToLoginInput(req *authv1.LoginRequest) authOperations.LoginInput {
	return authOperations.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func ToLogoutInput(req *authv1.LogoutRequest) authOperations.LogoutInput {
	return authOperations.LogoutInput{RefreshToken: req.GetRefreshToken()}
}

func ToRefreshInput(req *authv1.RefreshRequest) authOperations.RefreshInput {
	return authOperations.RefreshInput{RefreshToken: req.GetRefreshToken()}
}
