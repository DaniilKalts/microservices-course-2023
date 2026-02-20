package auth

import (
	"context"

	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
)

func (h *Handler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	tokens, err := h.authService.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &authv1.RefreshResponse{Tokens: toProtoTokenPair(tokens)}, nil
}
