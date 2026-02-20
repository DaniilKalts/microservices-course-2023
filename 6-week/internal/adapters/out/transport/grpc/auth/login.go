package auth

import (
	"context"

	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
)

func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokens, err := h.authService.Login(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &authv1.LoginResponse{Tokens: toProtoTokenPair(tokens)}, nil
}
