package mapper

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/6-week/internal/domain/user"
)

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

func ToProtoUsers(users []domainUser.User) []*userv1.User {
	protoUsers := make([]*userv1.User, 0, len(users))
	for i := range users {
		protoUsers = append(protoUsers, ToProtoUser(&users[i]))
	}

	return protoUsers
}
