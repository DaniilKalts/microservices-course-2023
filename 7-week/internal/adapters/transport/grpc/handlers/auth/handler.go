package auth

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	authService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/auth"
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

func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, authService.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, authService.ErrInvalidCredentials.Error())
	case errors.Is(err, authService.ErrInvalidRefreshToken):
		return status.Error(codes.Unauthenticated, authService.ErrInvalidRefreshToken.Error())
	case errors.Is(err, authService.ErrAuthentication):
		return status.Error(codes.Internal, authService.ErrAuthentication.Error())
	case errors.Is(err, authService.ErrUserIDEmpty):
		return status.Error(codes.Internal, authService.ErrUserIDEmpty.Error())
	case errors.Is(err, authService.ErrIssueTokens):
		return status.Error(codes.Internal, authService.ErrIssueTokens.Error())
	case errors.Is(err, domainUser.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, domainUser.ErrEmailAlreadyExists.Error())
	default:
		h.logger.Error("unhandled auth error", zap.Error(err))
		return status.Error(codes.Internal, "internal error")
	}
}
