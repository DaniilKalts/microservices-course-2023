package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/protoutil"
)

// Proto → Domain

func toCreateUserInput(req *userv1.CreateRequest) domainUser.CreateInput {
	return domainUser.CreateInput{
		User: &domainUser.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  domainUser.RoleUser,
		},
		Password: req.GetPassword(),
	}
}

func toUpdateUserInput(req *userv1.UpdateRequest) domainUser.UpdateInput {
	return domainUser.UpdateInput{
		ID:    req.GetId(),
		Name:  protoutil.StringValuePtr(req.GetName()),
		Email: protoutil.StringValuePtr(req.GetEmail()),
	}
}

// Domain → Proto

func toProtoUser(user *domainUser.User) *userv1.User {
	return &userv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      userv1.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func toProtoUsers(users []domainUser.User) []*userv1.User {
	protoUsers := make([]*userv1.User, 0, len(users))
	for i := range users {
		protoUsers = append(protoUsers, toProtoUser(&users[i]))
	}
	return protoUsers
}
