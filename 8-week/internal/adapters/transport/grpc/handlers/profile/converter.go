package profile

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/protoutil"
)

// Proto → Domain

func toUpdateProfileInput(userID string, req *userv1.UpdateProfileRequest) domainUser.UpdateInput {
	return domainUser.UpdateInput{
		ID:       userID,
		Name:     protoutil.StringValuePtr(req.GetName()),
		Email:    protoutil.StringValuePtr(req.GetEmail()),
		Password: protoutil.StringValuePtr(req.GetPassword()),
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
