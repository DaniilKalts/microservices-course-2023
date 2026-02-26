package grpc

import (
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	authHandler "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/auth"
	userHandler "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/interceptor"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

type Deps struct {
	Config     ServerConfig
	JWTManager jwt.Manager
	Logger     *zap.Logger
	Services   service.Services
}

func NewServer(deps Deps) (*grpc.Server, error) {
	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.LoggingInterceptor(deps.Logger),
			interceptor.AuthInterceptor(deps.JWTManager),
			interceptor.ValidationInterceptor(),
		),
	}

	cfg := deps.Config

	if cfg.EnableTLS {
		if cfg.CertFile == "" || cfg.KeyFile == "" {
			return nil, errors.New("grpc tls is enabled but cert or key file is empty")
		}

		creds, err := credentials.NewServerTLSFromFile(cfg.CertFile, cfg.KeyFile)
		if err != nil {
			return nil, err
		}

		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	server := grpc.NewServer(grpcOpts...)

	userv1.RegisterUserV1Server(server, userHandler.NewHandler(deps.Services.User))
	authv1.RegisterAuthV1Server(server, authHandler.NewHandler(deps.Services.Auth))

	reflection.Register(server)

	return server, nil
}
