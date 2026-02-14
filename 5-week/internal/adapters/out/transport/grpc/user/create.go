package user

import (
	"context"

	userv1 "github.com/DaniilKalts/microservices-course-2023/5-week/gen/go/user/v1"
)

func (h *Handler) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	userID, err := h.userService.Create(ctx, toDomainFromCreate(req), req.GetPassword(), req.GetPasswordConfirm())
	if err != nil {
		return nil, err
	}

	return &userv1.CreateResponse{Id: userID}, nil
}
