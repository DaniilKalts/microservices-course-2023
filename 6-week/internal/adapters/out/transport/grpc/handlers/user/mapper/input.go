package mapper

import (
	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/handlers/user/procedures"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

func ToCreateInput(req *userv1.CreateRequest) procedures.CreateInput {
	return procedures.CreateInput{
		User: &domainUser.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  domainUser.RoleUser,
		},
		Password:        req.GetPassword(),
		PasswordConfirm: req.GetPasswordConfirm(),
	}
}

func ToGetInput(req *userv1.GetRequest) procedures.GetInput {
	return procedures.GetInput{ID: req.GetId()}
}

func ToUpdateInput(req *userv1.UpdateRequest) procedures.UpdateInput {
	patch := &domainUser.UpdatePatch{}
	if req.GetName() != nil {
		name := req.GetName().GetValue()
		patch.Name = &name
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		patch.Email = &email
	}

	return procedures.UpdateInput{ID: req.GetId(), Patch: patch}
}

func ToDeleteInput(req *userv1.DeleteRequest) procedures.DeleteInput {
	return procedures.DeleteInput{ID: req.GetId()}
}
