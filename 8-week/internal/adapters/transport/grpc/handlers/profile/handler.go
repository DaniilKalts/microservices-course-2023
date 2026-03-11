package profile

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	userHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/interceptor/auth"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/protoutil"
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

	return &userv1.GetProfileResponse{User: userHandler.ToProtoUser(user)}, nil
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

func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, domainUser.ErrNotFound):
		return status.Error(codes.NotFound, domainUser.ErrNotFound.Error())
	case errors.Is(err, domainUser.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, domainUser.ErrEmailAlreadyExists.Error())
	case errors.Is(err, domainUser.ErrNoFieldsToUpdate):
		return status.Error(codes.InvalidArgument, domainUser.ErrNoFieldsToUpdate.Error())
	default:
		h.logger.Error("unhandled profile error", zap.Error(err))
		return status.Error(codes.Internal, "internal error")
	}
}

func toUpdateProfileInput(userID string, req *userv1.UpdateProfileRequest) userService.UpdateInput {
	return userService.UpdateInput{
		ID:       userID,
		Name:     protoutil.StringValuePtr(req.GetName()),
		Email:    protoutil.StringValuePtr(req.GetEmail()),
		Password: protoutil.StringValuePtr(req.GetPassword()),
	}
}
