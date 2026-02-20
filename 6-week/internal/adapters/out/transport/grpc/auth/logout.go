package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
)

func (h *Handler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.authService.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
