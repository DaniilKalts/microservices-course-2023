package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
)

type Handlers struct {
	User userv1.UserV1Server
	Auth authv1.AuthV1Server
}

func RegisterServices(server *grpc.Server, handlers Handlers) {
	userv1.RegisterUserV1Server(server, handlers.User)
	authv1.RegisterAuthV1Server(server, handlers.Auth)

	reflection.Register(server)
}
