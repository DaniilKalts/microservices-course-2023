package handlers

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/transport/grpc/interceptor"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/protoutil"
)

type ProfileHandler struct {
	userv1.UnimplementedProfileV1Server
	userService userService.Service
	logger      *zap.Logger
}

func NewProfileHandler(userService userService.Service, logger *zap.Logger) *ProfileHandler {
	return &ProfileHandler{userService: userService, logger: logger}
}

func (h *ProfileHandler) GetProfile(ctx context.Context, _ *emptypb.Empty) (*userv1.GetProfileResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.Get(ctx, userID)
	if err != nil {
		return nil, h.mapProfileError(err)
	}

	return &userv1.GetProfileResponse{User: toProtoUser(user)}, nil
}

func (h *ProfileHandler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Update(ctx, toUpdateProfileInput(userID, req)); err != nil {
		return nil, h.mapProfileError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ProfileHandler) DeleteProfile(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Delete(ctx, userID); err != nil {
		return nil, h.mapProfileError(err)
	}

	return &emptypb.Empty{}, nil
}

func currentUserID(ctx context.Context) (string, error) {
	claims, ok := interceptor.ClaimsFromContext(ctx)
	if !ok || claims == nil || claims.UserID == "" {
		return "", status.Error(codes.Unauthenticated, interceptor.ErrInvalidAccessToken.Error())
	}

	return claims.UserID, nil
}

func (h *ProfileHandler) mapProfileError(err error) error {
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
