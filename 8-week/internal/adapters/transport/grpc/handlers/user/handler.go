package user

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	userService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/user"
)

type Handler struct {
	userv1.UnimplementedUserV1Server
	userService userService.Service
	logger      *zap.Logger
}

func NewHandler(userService userService.Service, logger *zap.Logger) *Handler {
	return &Handler{userService: userService, logger: logger}
}

func (h *Handler) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	userID, err := h.userService.Create(ctx, toCreateUserInput(req))
	if err != nil {
		return nil, h.mapError(err)
	}

	return &userv1.CreateResponse{Id: userID}, nil
}

func (h *Handler) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	user, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, h.mapError(err)
	}

	return &userv1.GetResponse{User: toProtoUser(user)}, nil
}

func (h *Handler) List(ctx context.Context, _ *emptypb.Empty) (*userv1.ListResponse, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		return nil, h.mapError(err)
	}

	return &userv1.ListResponse{Users: toProtoUsers(users)}, nil
}

func (h *Handler) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	if err := h.userService.Update(ctx, toUpdateUserInput(req)); err != nil {
		return nil, h.mapError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, h.mapError(err)
	}

	return &emptypb.Empty{}, nil
}
