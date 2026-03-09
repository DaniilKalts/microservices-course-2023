package handlers

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/protoutil"
	userService "github.com/DaniilKalts/microservices-course-2023/7-week/internal/service/user"
)

type UserHandler struct {
	userv1.UnimplementedUserV1Server
	userService userService.Service
	logger      *zap.Logger
}

func NewUserHandler(userService userService.Service, logger *zap.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

func (h *UserHandler) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	userID, err := h.userService.Create(ctx, toCreateUserInput(req))
	if err != nil {
		return nil, h.mapUserError(err)
	}

	return &userv1.CreateResponse{Id: userID}, nil
}

func (h *UserHandler) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	user, err := h.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, h.mapUserError(err)
	}

	return &userv1.GetResponse{User: toProtoUser(user)}, nil
}

func (h *UserHandler) List(ctx context.Context, _ *emptypb.Empty) (*userv1.ListResponse, error) {
	users, err := h.userService.List(ctx)
	if err != nil {
		return nil, h.mapUserError(err)
	}

	return &userv1.ListResponse{Users: toProtoUsers(users)}, nil
}

func (h *UserHandler) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	if err := h.userService.Update(ctx, toUpdateUserInput(req)); err != nil {
		return nil, h.mapUserError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	if err := h.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, h.mapUserError(err)
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) mapUserError(err error) error {
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

func toProtoUser(user *domainUser.User) *userv1.User {
	return &userv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      userv1.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func toProtoUsers(users []domainUser.User) []*userv1.User {
	protoUsers := make([]*userv1.User, 0, len(users))
	for i := range users {
		protoUsers = append(protoUsers, toProtoUser(&users[i]))
	}

	return protoUsers
}

func toCreateUserInput(req *userv1.CreateRequest) userService.CreateInput {
	return userService.CreateInput{
		User: &domainUser.User{
			Name:  req.GetName(),
			Email: req.GetEmail(),
			Role:  domainUser.RoleUser,
		},
		Password: req.GetPassword(),
	}
}

func toUpdateUserInput(req *userv1.UpdateRequest) userService.UpdateInput {
	return userService.UpdateInput{
		ID:    req.GetId(),
		Name:  protoutil.StringValuePtr(req.GetName()),
		Email: protoutil.StringValuePtr(req.GetEmail()),
	}
}
