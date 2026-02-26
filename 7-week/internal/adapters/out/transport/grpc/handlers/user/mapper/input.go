package mapper

import (
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user/operations"
)

func ToCreateInput(req *userv1.CreateRequest) userOperations.CreateInput {
	return userOperations.CreateInput{
		User: &domainUser.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  domainUser.RoleUser,
		},
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func ToGetInput(req *userv1.GetRequest) userOperations.GetInput {
	return userOperations.GetInput{ID: req.GetId()}
}

func ToUpdateInput(req *userv1.UpdateRequest) userOperations.UpdateInput {
	input := userOperations.UpdateInput{ID: req.GetId()}
	if req.GetName() != nil {
		name := req.GetName().GetValue()
		input.Name = &name
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		input.Email = &email
	}

	return input
}

func ToDeleteInput(req *userv1.DeleteRequest) userOperations.DeleteInput {
	return userOperations.DeleteInput{ID: req.GetId()}
}
