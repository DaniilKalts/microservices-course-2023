package grpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	authHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/auth"
	profileHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/profile"
	userHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/interceptor"
	authInterceptor "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/interceptor/auth"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string

	RequestTimeout time.Duration
}

type Deps struct {
	Config     ServerConfig
	Logger     *zap.Logger
	JWTManager jwt.Manager

	Services   service.Services
	AuthPolicy authInterceptor.AccessPolicy

	Tracer   opentracing.Tracer
	Registry prometheus.Registerer
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

	authIntcpt, err := authInterceptor.NewInterceptor(deps.JWTManager, authPolicy, logger.Named("interceptor.auth"))
	if err != nil {
		return nil, fmt.Errorf("init auth interceptor: %w", err)
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			interceptor.TimeoutInterceptor(deps.Config.RequestTimeout),
			interceptor.MetricsInterceptor(deps.Registry),
			interceptor.TracingInterceptor(deps.Tracer),
			interceptor.LoggingInterceptor(logger.Named("interceptor.logging")),
			authIntcpt,
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
			return nil, fmt.Errorf("load grpc tls credentials: %w", err)
		}

		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	server := grpc.NewServer(grpcOpts...)

	handlerLogger := logger.Named("handler")
	userv1.RegisterUserV1Server(server, userHandler.NewHandler(deps.Services.User, handlerLogger.Named("user")))
	userv1.RegisterProfileV1Server(server, profileHandler.NewHandler(deps.Services.User, handlerLogger.Named("profile")))
	authv1.RegisterAuthV1Server(server, authHandler.NewHandler(deps.Services.Auth, handlerLogger.Named("auth")))

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(server)

	return server, nil
}
