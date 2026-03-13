package auth

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/auth/v1"
	authService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/auth"
)

type Handler struct {
	authv1.UnimplementedAuthV1Server
	authService authService.Service
	logger      *zap.Logger
}

func NewHandler(authService authService.Service, logger *zap.Logger) *Handler {
	return &Handler{authService: authService, logger: logger}
}

func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, tokens, err := h.authService.Register(ctx, toRegisterInput(req))
	if err != nil {
		return nil, h.mapError(err)
	}

	return &authv1.RegisterResponse{User: toProtoRegisterUser(user), Tokens: toProtoTokenPair(tokens)}, nil
}

func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokens, err := h.authService.Login(ctx, toLoginInput(req))
	if err != nil {
		return nil, h.mapError(err)
	}

	return &authv1.LoginResponse{Tokens: toProtoTokenPair(tokens)}, nil
}

func (h *Handler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.authService.Logout(ctx, toLogoutInput(req)); err != nil {
		return nil, h.mapError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	tokens, err := h.authService.Refresh(ctx, toRefreshInput(req))
	if err != nil {
		return nil, h.mapError(err)
	}

	return &authv1.RefreshResponse{Tokens: toProtoTokenPair(tokens)}, nil
}
