package user

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	mapper "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/handlers/user/mapper"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/handlers/user/procedures"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
)

type Handler struct {
	userv1.UnimplementedUserV1Server
	userService service.UserService
}

func NewHandler(userService service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) Create(ctx context.Context, req *userv1.CreateRequest) (*userv1.CreateResponse, error) {
	userID, err := procedures.Create(ctx, h.userService, mapper.ToCreateInput(req))
	if err != nil {
		return nil, err
	}

	return &userv1.CreateResponse{Id: userID}, nil
}

func (h *Handler) Get(ctx context.Context, req *userv1.GetRequest) (*userv1.GetResponse, error) {
	entity, err := procedures.Get(ctx, h.userService, mapper.ToGetInput(req))
	if err != nil {
		return nil, err
	}

	return &userv1.GetResponse{User: mapper.ToProtoUser(entity)}, nil
}

func (h *Handler) List(ctx context.Context, _ *emptypb.Empty) (*userv1.ListResponse, error) {
	users, err := procedures.List(ctx, h.userService)
	if err != nil {
		return nil, err
	}

	return &userv1.ListResponse{Users: mapper.ToProtoUsers(users)}, nil
}

func (h *Handler) Update(ctx context.Context, req *userv1.UpdateRequest) (*emptypb.Empty, error) {
	if err := procedures.Update(ctx, h.userService, mapper.ToUpdateInput(req)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *Handler) Delete(ctx context.Context, req *userv1.DeleteRequest) (*emptypb.Empty, error) {
	if err := procedures.Delete(ctx, h.userService, mapper.ToDeleteInput(req)); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
