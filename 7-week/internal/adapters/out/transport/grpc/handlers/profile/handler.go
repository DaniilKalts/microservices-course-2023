package profile

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	profileMapper "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/profile/mapper"
	userMapper "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/user/mapper"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/interceptor"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	userOperations "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user/operations"
)

type Handler struct {
	userv1.UnimplementedProfileV1Server
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) GetProfile(ctx context.Context, _ *emptypb.Empty) (*userv1.GetProfileResponse, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	entity, err := h.userService.Get(ctx, userOperations.GetInput{ID: userID})
	if err != nil {
		return nil, err
	}

	return &userv1.GetProfileResponse{User: userMapper.ToProtoUser(entity)}, nil
}

func (h *Handler) UpdateProfile(ctx context.Context, req *userv1.UpdateProfileRequest) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Update(ctx, profileMapper.ToUpdateInput(userID, req)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteProfile(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}

	if err = h.userService.Delete(ctx, userOperations.DeleteInput{ID: userID}); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func currentUserID(ctx context.Context) (string, error) {
	claims, ok := interceptor.ClaimsFromContext(ctx)
	if !ok || claims == nil || claims.UserID == "" {
		return "", status.Error(codes.Unauthenticated, "invalid access token")
	}

	return claims.UserID, nil
}
