package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/5-week/gen/go/user/v1"
)

func (h *Handler) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	if err := h.userService.Update(ctx, req.GetId(), toDomainPatchFromUpdate(req)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
