package profile

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/interceptor/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user"
)

type Handler struct {
	userv1.UnimplementedProfileV1Server
	userService userService.Service
	logger      *zap.Logger
}

func NewHandler(userService userService.Service, logger *zap.Logger) *Handler {
	return &Handler{userService: userService, logger: logger}
}

func (h *Handler) GetProfile(ctx context.Context, _ *emptypb.Empty) (*userv1.GetProfileResponse, error) {
	userID, err := h.getIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.Get(ctx, userID)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &userv1.GetProfileResponse{User: toProtoUser(user)}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*emptypb.Empty, error) {
	userID, err := h.getIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Update(ctx, toUpdateProfileInput(userID, req)); err != nil {
		return nil, h.mapError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteProfile(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userID, err := h.getIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Delete(ctx, userID); err != nil {
		return nil, h.mapError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) getIDFromContext(ctx context.Context) (string, error) {
	claims, ok := auth.ClaimsFromContext(ctx)
	if !ok || claims == nil || claims.UserID == "" {
		return "", status.Error(codes.Unauthenticated, auth.ErrInvalidAccessToken.Error())
	}

	return claims.UserID, nil
}
