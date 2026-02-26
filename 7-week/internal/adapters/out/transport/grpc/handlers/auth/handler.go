package auth

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	mapper "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/auth/mapper"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
)

type Handler struct {
	authv1.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewHandler(authService service.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, tokens, err := h.authService.Register(ctx, mapper.ToRegisterInput(req))
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &authv1.RegisterResponse{User: mapper.ToProtoRegisterUser(user), Tokens: mapper.ToProtoTokenPair(tokens)}, nil
}

func (h *Handler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokens, err := h.authService.Login(ctx, mapper.ToLoginInput(req))
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &authv1.LoginResponse{Tokens: mapper.ToProtoTokenPair(tokens)}, nil
}

func (h *Handler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.authService.Logout(ctx, mapper.ToLogoutInput(req)); err != nil {
		return nil, mapAuthError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	tokens, err := h.authService.Refresh(ctx, mapper.ToRefreshInput(req))
	if err != nil {
		return nil, mapAuthError(err)
	}

	return &authv1.RefreshResponse{Tokens: mapper.ToProtoTokenPair(tokens)}, nil
}
