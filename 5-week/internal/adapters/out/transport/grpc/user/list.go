package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/5-week/gen/go/user/v1"
)

func (h *Handler) List(ctx context.Context, _ *emptypb.Empty) (*userv1.ListResponse, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		return nil, err
	}

	return &userv1.ListResponse{Users: toProtoUsers(users)}, nil
}
