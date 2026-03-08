package profile

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	profileMapper "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/grpc/handlers/profile/mapper"
	userMapper "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/grpc/handlers/user/mapper"
	authInterceptor "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/grpc/interceptor/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
)

type Handler struct {
	userv1.UnimplementedProfileV1Server
	userService userService.Service
}

func NewHandler(userService userService.Service) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) GetProfile(ctx context.Context, _ *emptypb.Empty) (*userv1.GetProfileResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := h.userService.Get(ctx, userID)
	if err != nil {
		return nil, mapProfileError(err)
	}

	return &userv1.GetProfileResponse{User: userMapper.ToProtoUser(entity)}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Update(ctx, profileMapper.ToUpdateInput(userID, req)); err != nil {
		return nil, mapProfileError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteProfile(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Delete(ctx, userID); err != nil {
		return nil, mapProfileError(err)
	}

	return &emptypb.Empty{}, nil
}

func currentUserID(ctx context.Context) (string, error) {
	claims, ok := authInterceptor.ClaimsFromContext(ctx)
	if !ok || claims == nil || claims.UserID == "" {
		return "", status.Error(codes.Unauthenticated, authInterceptor.ErrInvalidAccessToken.Error())
	}

	return claims.UserID, nil
}
