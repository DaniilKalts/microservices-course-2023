package mapper

import (
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/protoutil"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
)

func ToUpdateInput(userID string, req *userv1.UpdateProfileRequest) userService.UpdateInput {
	return userService.UpdateInput{
		ID:       userID,
		Name:     protoutil.StringValuePtr(req.GetName()),
		Email:    protoutil.StringValuePtr(req.GetEmail()),
		Password: protoutil.StringValuePtr(req.GetPassword()),
	}
}
