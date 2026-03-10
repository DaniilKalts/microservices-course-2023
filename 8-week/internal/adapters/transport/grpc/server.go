package grpc

import (
	"errors"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/gen/grpc/user/v1"
	authHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/auth"
	profileHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/profile"
	userHandler "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/handlers/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc/interceptor"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
)

type ServerConfig struct {
	EnableTLS bool
	CertFile  string
	KeyFile   string
}

type Deps struct {
	Config     ServerConfig
	Logger     *zap.Logger
	JWTManager jwt.Manager

	Services   service.Services
	AuthPolicy interceptor.AccessPolicy

	Tracer opentracing.Tracer
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
			interceptor.MetricsInterceptor(),
			interceptor.TracingInterceptor(deps.Tracer),
			interceptor.LoggingInterceptor(logger.Named("interceptor.logging")),
			interceptor.AuthInterceptor(deps.JWTManager, authPolicy, logger.Named("interceptor.auth")),
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

	reflection.Register(server)

	return server, nil
}

func defaultAccessPolicy() (interceptor.AccessPolicy, error) {
	public := interceptor.PublicGroup()
	authenticated := interceptor.RoleGroup("authenticated", int32(domainUser.RoleUser), int32(domainUser.RoleAdmin))
	admin := interceptor.RoleGroup("admin", int32(domainUser.RoleAdmin))

	return interceptor.NewAccessPolicy(
		interceptor.MethodGroup{
			Group: public,
			Methods: []string{
				authv1.AuthV1_Register_FullMethodName,
				authv1.AuthV1_Login_FullMethodName,
				authv1.AuthV1_Refresh_FullMethodName,
				userv1.UserV1_List_FullMethodName,
				userv1.UserV1_Get_FullMethodName,
			},
		},
		interceptor.MethodGroup{
			Group: admin,
			Methods: []string{
				userv1.UserV1_Create_FullMethodName,
				userv1.UserV1_Update_FullMethodName,
				userv1.UserV1_Delete_FullMethodName,
			},
		},
		interceptor.MethodGroup{
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
