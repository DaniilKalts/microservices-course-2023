package grpc

import (
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	authHandler "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/auth"
	profileHandler "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/profile"
	userHandler "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/handlers/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/interceptor"
	authInterceptor "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc/interceptor/auth"
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
	AuthPolicy authInterceptor.AccessPolicy
	Tracer     opentracing.Tracer
}

func NewServer(deps Deps) (*grpc.Server, error) {
	logger := deps.Logger
	if logger == nil {
		return nil, errors.New("grpc logger is nil")
	}

	authPolicy := deps.AuthPolicy
	if authPolicy.IsEmpty() {
		var err error
		authPolicy, err = authInterceptor.DefaultAccessPolicy()
		if err != nil {
			return nil, fmt.Errorf("build auth access policy: %w", err)
		}
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.MetricsInterceptor(),
			interceptor.TracingInterceptor(deps.Tracer),
			interceptor.LoggingInterceptor(logger.Named("interceptor.logging")),
			authInterceptor.AuthInterceptor(deps.JWTManager, authPolicy),
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
	userv1.RegisterProfileV1Server(server, profileHandler.NewHandler(deps.Services.User))
	authv1.RegisterAuthV1Server(server, authHandler.NewHandler(deps.Services.Auth))

	reflection.Register(server)

	return server, nil
}
