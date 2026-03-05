package mapper

import (
	"google.golang.org/protobuf/types/known/wrapperspb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	userOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user/operations"
)

func ToUpdateInput(userID string, req *userv1.UpdateProfileRequest) userOperations.UpdateInput {
	return userOperations.UpdateInput{
		ID:       userID,
		Name:     stringValuePtr(req.GetName()),
		Email:    stringValuePtr(req.GetEmail()),
		Password: stringValuePtr(req.GetPassword()),
	}
}

func stringValuePtr(value *wrapperspb.StringValue) *string {
	if value == nil {
		return nil
	}

	str := value.GetValue()
	if str == "" {
		return nil
	}

	return &str
}
