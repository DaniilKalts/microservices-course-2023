package grpc

import (
	"errors"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/service"
	authService "github.com/DaniilKalts/microservices-course-2023/8-week/internal/service/auth"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string

	RequestTimeout time.Duration

	RateLimitRPS   float64
	RateLimitBurst int

	RateLimitAuthRPS   float64
	RateLimitAuthBurst int
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

var defaultErrorRules = []interceptor.ErrorRule{
	// Domain errors
	{
		Target: domainUser.ErrNotFound,
		Code:   codes.NotFound,
	},
	{
		Target: domainUser.ErrEmailAlreadyExists,
		Code:   codes.AlreadyExists,
	},
	{
		Target: domainUser.ErrNoFieldsToUpdate,
		Code:   codes.InvalidArgument,
	},
	{
		Target: domainUser.ErrWeakPassword,
		Code:   codes.InvalidArgument,
	},

	// Auth service errors
	{
		Target: authService.ErrInvalidCredentials,
		Code:   codes.Unauthenticated,
	},
	{
		Target: authService.ErrInvalidRefreshToken,
		Code:   codes.Unauthenticated,
	},
	{
		Target:  authService.ErrAuthentication,
		Code:    codes.Internal,
		Message: "authentication failed",
	},
	{
		Target:  authService.ErrUserIDEmpty,
		Code:    codes.Internal,
		Message: "internal error",
	},
	{
		Target:  authService.ErrIssueTokens,
		Code:    codes.Internal,
		Message: "internal error",
	},
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
			interceptor.RateLimitInterceptor(interceptor.RateLimitConfig{
				RPS:       deps.Config.RateLimitRPS,
				Burst:     deps.Config.RateLimitBurst,
				AuthRPS:   deps.Config.RateLimitAuthRPS,
				AuthBurst: deps.Config.RateLimitAuthBurst,
			}),
			interceptor.TimeoutInterceptor(deps.Config.RequestTimeout),
			interceptor.MetricsInterceptor(deps.Registry),
			interceptor.TracingInterceptor(deps.Tracer),
			interceptor.LoggingInterceptor(logger.Named("interceptor.logging")),
			authIntcpt,
			interceptor.ValidationInterceptor(),
			interceptor.ErrorMappingInterceptor(logger.Named("interceptor.errors"), defaultErrorRules),
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

	userv1.RegisterUserV1Server(server, userHandler.NewHandler(deps.Services.User))
	userv1.RegisterProfileV1Server(server, profileHandler.NewHandler(deps.Services.User))
	authv1.RegisterAuthV1Server(server, authHandler.NewHandler(deps.Services.Auth))

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(server, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(server)

	return server, nil
}
