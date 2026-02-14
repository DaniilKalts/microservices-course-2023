package user

import (
	"context"

	userv1 "github.com/DaniilKalts/microservices-course-2023/5-week/gen/go/user/v1"
)

func (h *Handler) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	entity, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, err
	}

	return &userv1.GetResponse{User: toProtoUser(entity)}, nil
}
