package auth

import (
	"context"

	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
)

func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	userID, tokens, err := h.authService.Register(ctx, req.GetName(), req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, err
	}

	return &authv1.RegisterResponse{Id: userID, Tokens: toProtoTokenPair(tokens)}, nil
}
