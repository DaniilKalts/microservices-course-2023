package handlers

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

type AuthHandler struct {
	authv1.UnimplementedAuthV1Server
	authService authService.Service
	logger      *zap.Logger
}

func NewAuthHandler(authService authService.Service, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{authService: authService, logger: logger}
}

func (h *AuthHandler) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, tokens, err := h.authService.Register(ctx, toRegisterInput(req))
	if err != nil {
		return nil, h.mapAuthError(err)
	}

	return &authv1.RegisterResponse{User: toProtoRegisterUser(user), Tokens: toProtoTokenPair(tokens)}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokens, err := h.authService.Login(ctx, toLoginInput(req))
	if err != nil {
		return nil, h.mapAuthError(err)
	}

	return &authv1.LoginResponse{Tokens: toProtoTokenPair(tokens)}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*emptypb.Empty, error) {
	if err := h.authService.Logout(ctx, toLogoutInput(req)); err != nil {
		return nil, h.mapAuthError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *AuthHandler) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	tokens, err := h.authService.Refresh(ctx, toRefreshInput(req))
	if err != nil {
		return nil, h.mapAuthError(err)
	}

	return &authv1.RefreshResponse{Tokens: toProtoTokenPair(tokens)}, nil
}

func (h *AuthHandler) mapAuthError(err error) error {
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

func toRegisterInput(req *authv1.RegisterRequest) authService.RegisterInput {
	return authService.RegisterInput{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func toLoginInput(req *authv1.LoginRequest) authService.LoginInput {
	return authService.LoginInput{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}
}

func toLogoutInput(req *authv1.LogoutRequest) authService.LogoutInput {
	return authService.LogoutInput{RefreshToken: req.GetRefreshToken()}
}

func toRefreshInput(req *authv1.RefreshRequest) authService.RefreshInput {
	return authService.RefreshInput{RefreshToken: req.GetRefreshToken()}
}

func toProtoTokenPair(tokens authService.TokenPair) *authv1.TokenPair {
	return &authv1.TokenPair{
		AccessToken:           tokens.AccessToken,
		RefreshToken:          tokens.RefreshToken,
		AccessTokenExpiresIn:  tokens.AccessTokenExpiresIn,
		RefreshTokenExpiresIn: tokens.RefreshTokenExpiresIn,
	}
}

func toProtoRegisterUser(user domainUser.User) *authv1.RegisterResponse_User {
	return &authv1.RegisterResponse_User{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}
