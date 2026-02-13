package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/4-week/gen/go/user/v1"
)

func (h *Handler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
