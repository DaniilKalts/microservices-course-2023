package mapper

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/handlers/auth/procedures"
)

func ToRegisterInput(req *authv1.RegisterRequest) procedures.RegisterInput {
	return procedures.RegisterInput{
		Name:            req.GetName(),
		Email:           req.GetEmail(),
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func ToLoginInput(req *authv1.LoginRequest) procedures.LoginInput {
	return procedures.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func ToLogoutInput(req *authv1.LogoutRequest) procedures.LogoutInput {
	return procedures.LogoutInput{RefreshToken: req.GetRefreshToken()}
}

func ToRefreshInput(req *authv1.RefreshRequest) procedures.RefreshInput {
	return procedures.RefreshInput{RefreshToken: req.GetRefreshToken()}
}
