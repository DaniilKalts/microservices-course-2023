package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/protoutil"
)

// ToProtoUser converts a domain user to a proto user.
// Exported for use by the profile handler.
func ToProtoUser(user *domainUser.User) *userv1.User {
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
		protoUsers = append(protoUsers, ToProtoUser(&users[i]))
	}

	return protoUsers
}

func toCreateUserInput(req *userv1.CreateRequest) userService.CreateInput {
	return userService.CreateInput{
		User: &domainUser.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  domainUser.RoleUser,
		},
		Password: req.GetPassword(),
	}
}

func toUpdateUserInput(req *userv1.UpdateRequest) userService.UpdateInput {
	return userService.UpdateInput{
		ID:    req.GetId(),
		Name:  protoutil.StringValuePtr(req.GetName()),
		Email: protoutil.StringValuePtr(req.GetEmail()),
	}
}
