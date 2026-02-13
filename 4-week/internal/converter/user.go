package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/4-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/4-week/internal/models"
)

func ToUserFromService(user *models.User) *userv1.User {
	var createdAt *timestamppb.Timestamp = timestamppb.New(user.CreatedAt)
	var updatedAt *timestamppb.Timestamp = timestamppb.New(*user.UpdatedAt)

	return &userv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      userv1.Role(user.Role),
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
