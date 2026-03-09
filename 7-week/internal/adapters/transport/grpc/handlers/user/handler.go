package user

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
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

	return &userv1.GetResponse{User: ToProtoUser(user)}, nil
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

func (h *Handler) mapError(err error) error {
	switch {
	case errors.Is(err, domainUser.ErrNotFound):
		return status.Error(codes.NotFound, domainUser.ErrNotFound.Error())
	case errors.Is(err, domainUser.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, domainUser.ErrEmailAlreadyExists.Error())
	case errors.Is(err, domainUser.ErrNoFieldsToUpdate):
		return status.Error(codes.InvalidArgument, domainUser.ErrNoFieldsToUpdate.Error())
	default:
		h.logger.Error("unhandled user error", zap.Error(err))
		return status.Error(codes.Internal, "internal error")
	}
}
