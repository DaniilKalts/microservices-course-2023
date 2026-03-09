package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/transport/grpc/middleware"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/protoutil"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
)

type ProfileHandler struct {
	userv1.UnimplementedProfileV1Server
	userService userService.Service
}

func NewProfileHandler(userService userService.Service) *ProfileHandler {
	return &ProfileHandler{userService: userService}
}

func (h *ProfileHandler) GetProfile(ctx context.Context, _ *emptypb.Empty) (*userv1.GetProfileResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	user, err := h.userService.Get(ctx, userID)
	if err != nil {
		return nil, mapDomainUserError(err)
	}

	return &userv1.GetProfileResponse{User: toProtoUser(user)}, nil
}

func (h *ProfileHandler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Update(ctx, toUpdateProfileInput(userID, req)); err != nil {
		return nil, mapDomainUserError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *ProfileHandler) DeleteProfile(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Delete(ctx, userID); err != nil {
		return nil, mapDomainUserError(err)
	}

	return &emptypb.Empty{}, nil
}

func currentUserID(ctx context.Context) (string, error) {
	claims, ok := middleware.ClaimsFromContext(ctx)
	if !ok || claims == nil || claims.UserID == "" {
		return "", status.Error(codes.Unauthenticated, middleware.ErrInvalidAccessToken.Error())
	}

	return claims.UserID, nil
}

func toUpdateProfileInput(userID string, req *userv1.UpdateProfileRequest) userService.UpdateInput {
	return userService.UpdateInput{
		ID:       userID,
		Name:     protoutil.StringValuePtr(req.GetName()),
		Email:    protoutil.StringValuePtr(req.GetEmail()),
		Password: protoutil.StringValuePtr(req.GetPassword()),
	}
}
