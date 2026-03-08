package mapper

import (
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/protoutil"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
)

func ToCreateInput(req *userv1.CreateRequest) userService.CreateInput {
	return userService.CreateInput{
		User: &domainUser.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  domainUser.RoleUser,
		},
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func ToUpdateInput(req *userv1.UpdateRequest) userService.UpdateInput {
	return userService.UpdateInput{
		ID:    req.GetId(),
		Name:  protoutil.StringValuePtr(req.GetName()),
		Email: protoutil.StringValuePtr(req.GetEmail()),
	}
}
