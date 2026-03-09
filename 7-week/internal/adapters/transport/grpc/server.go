package grpc

import (
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/7-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/transport/grpc/handlers"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/transport/grpc/middleware"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

type Deps struct {
	Config ServerConfig
	Logger     *zap.Logger
	JWTManager jwt.Manager

	Services   service.Services
	AuthPolicy middleware.AccessPolicy

	Tracer          opentracing.Tracer
	ResponseCounter *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewServer(deps Deps) (*grpc.Server, error) {
	logger := deps.Logger
	if logger == nil {
		return nil, errors.New("grpc logger is nil")
	}

	authPolicy := deps.AuthPolicy
	if authPolicy.IsEmpty() {
		var err error
		authPolicy, err = defaultAccessPolicy()
		if err != nil {
			return nil, fmt.Errorf("build auth access policy: %w", err)
		}
	}

	grpcOpts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			middleware.MetricsInterceptor(deps.ResponseCounter, deps.RequestDuration),
			middleware.TracingInterceptor(deps.Tracer),
			middleware.LoggingInterceptor(logger.Named("interceptor.logging")),
			middleware.AuthInterceptor(deps.JWTManager, authPolicy),
			middleware.ValidationInterceptor(),
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

	userv1.RegisterUserV1Server(server, handlers.NewUserHandler(deps.Services.User))
	userv1.RegisterProfileV1Server(server, handlers.NewProfileHandler(deps.Services.User))
	authv1.RegisterAuthV1Server(server, handlers.NewAuthHandler(deps.Services.Auth))

	reflection.Register(server)

	return server, nil
}

func defaultAccessPolicy() (middleware.AccessPolicy, error) {
	public := middleware.PublicGroup()
	authenticated := middleware.RoleGroup("authenticated", int32(domainUser.RoleUser), int32(domainUser.RoleAdmin))
	admin := middleware.RoleGroup("admin", int32(domainUser.RoleAdmin))

	return middleware.NewAccessPolicy(
		middleware.MethodGroup{
			Group: public,
			Methods: []string{
				authv1.AuthV1_Register_FullMethodName,
				authv1.AuthV1_Login_FullMethodName,
				authv1.AuthV1_Refresh_FullMethodName,
				userv1.UserV1_List_FullMethodName,
				userv1.UserV1_Get_FullMethodName,
			},
		},
		middleware.MethodGroup{
			Group: admin,
			Methods: []string{
				userv1.UserV1_Create_FullMethodName,
				userv1.UserV1_Update_FullMethodName,
				userv1.UserV1_Delete_FullMethodName,
			},
		},
		middleware.MethodGroup{
			Group: authenticated,
			Methods: []string{
				authv1.AuthV1_Logout_FullMethodName,
				userv1.ProfileV1_GetProfile_FullMethodName,
				userv1.ProfileV1_UpdateProfile_FullMethodName,
				userv1.ProfileV1_DeleteProfile_FullMethodName,
			},
		},
	)
}
