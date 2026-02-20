package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	authv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/in/database/postgres"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/in/transport/http/swagger"
	authAPI "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/auth"
	grpcInterceptor "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/interceptor"
	userAPI "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	userRepository "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
	authService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/auth"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/pkg/jwt"
)

const (
	swaggerBasePath         = "/swagger"
	swaggerMergedOpenAPI    = "gen/openapi/gateway.swagger.json"
	swaggerUserBasePath     = swaggerBasePath + "/user"
	swaggerUserOpenAPI      = "gen/openapi/user/v1/user.swagger.json"
	swaggerAuthBasePath     = swaggerBasePath + "/auth"
	swaggerAuthOpenAPI      = "gen/openapi/auth/v1/auth.swagger.json"
	swaggerRedirectTarget   = swaggerBasePath + "/"
	swaggerUserRedirectPath = swaggerUserBasePath + "/"
	swaggerAuthRedirectPath = swaggerAuthBasePath + "/"
)

type Container struct {
	Cfg config.Config

	DB         database.Client
	Tx         database.TxManager
	JWTManager jwt.Manager

	UserRepo    repository.UserRepository
	UserSvc     service.UserService
	userHandler userv1.UserV1Server

	AuthSvc     service.AuthService
	authHandler authv1.AuthV1Server

	GRPC    *grpc.Server
	Gateway http.Handler
}

func Build(ctx context.Context, configPath string) (*Container, error) {
	container := &Container{}

	if err := container.initConfig(configPath); err != nil {
		return nil, err
	}
	if err := container.initDatabase(ctx); err != nil {
		return nil, err
	}
	if err := container.initJWTManager(); err != nil {
		return nil, err
	}

	container.initTxManager()
	container.initUserRepository()
	container.initUserService()
	container.initAuthService()
	container.initUserHandler()
	container.initAuthHandler()

	if err := container.initGRPC(); err != nil {
		return nil, err
	}

	if err := container.initGateway(ctx); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initConfig(configPath string) error {
	config.Load(configPath)

	cfg, err := env.NewConfig()
	if err != nil {
		return fmt.Errorf("load env config: %w", err)
	}

	c.Cfg = cfg

	return nil
}

func (c *Container) initDatabase(ctx context.Context) error {
	db, err := postgres.New(ctx, c.Cfg.Postgres().DSN())
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	c.DB = db

	return nil
}

func (c *Container) initTxManager() {
	c.Tx = transaction.NewTransactionManager(c.DB.DB())
}

func (c *Container) initUserRepository() {
	c.UserRepo = userRepository.NewRepository(c.DB)
}

func (c *Container) initJWTManager() error {
	privateKey, err := jwt.LoadPrivateKey(c.Cfg.JWT().PrivateKeyFile())
	if err != nil {
		return fmt.Errorf("load jwt private key: %w", err)
	}

	publicKey, err := jwt.LoadPublicKey(c.Cfg.JWT().PublicKeyFile())
	if err != nil {
		return fmt.Errorf("load jwt public key: %w", err)
	}

	jwtManager, err := jwt.NewManager(privateKey, publicKey, jwt.Config{
		Issuer:          c.Cfg.JWT().Issuer(),
		Subject:         c.Cfg.JWT().Subject(),
		Audience:        c.Cfg.JWT().Audience(),
		AccessTokenTTL:  c.Cfg.JWT().AccessExpiresAt(),
		RefreshTokenTTL: c.Cfg.JWT().RefreshExpiresAt(),
		NotBeforeOffset: c.Cfg.JWT().NotBefore(),
		IssuedAtOffset:  c.Cfg.JWT().IssuedAt(),
	})
	if err != nil {
		return fmt.Errorf("init jwt manager: %w", err)
	}

	c.JWTManager = jwtManager

	return nil
}

func (c *Container) initUserService() {
	c.UserSvc = userService.NewService(c.UserRepo)
}

func (c *Container) initAuthService() {
	c.AuthSvc = authService.NewService(c.UserSvc, c.JWTManager)
}

func (c *Container) initUserHandler() {
	c.userHandler = userAPI.NewHandler(c.UserSvc)
}

func (c *Container) initAuthHandler() {
	c.authHandler = authAPI.NewHandler(c.AuthSvc)
}

func (c *Container) initGRPC() error {
	grpcOpts := make([]grpc.ServerOption, 0, 2)
	grpcOpts = append(grpcOpts, grpc.ChainUnaryInterceptor(grpcInterceptor.ValidationInterceptor()))

	if c.Cfg.TLS().Enabled() {
		cert := c.Cfg.TLS().CertFile()
		key := c.Cfg.TLS().KeyFile()

		creds, err := credentials.NewServerTLSFromFile(cert, key)
		if err != nil {
			return err
		}

		grpcOpts = append(grpcOpts, grpc.Creds(creds))
	}

	grpcServer := grpc.NewServer(grpcOpts...)
	reflection.Register(grpcServer)

	userv1.RegisterUserV1Server(grpcServer, c.userHandler)
	authv1.RegisterAuthV1Server(grpcServer, c.authHandler)

	c.GRPC = grpcServer

	return nil
}

func (c *Container) initGateway(ctx context.Context) error {
	gatewayMux := runtime.NewServeMux()
	if err := userv1.RegisterUserV1HandlerServer(ctx, gatewayMux, c.userHandler); err != nil {
		return fmt.Errorf("register grpc-gateway handlers: %w", err)
	}
	if err := authv1.RegisterAuthV1HandlerServer(ctx, gatewayMux, c.authHandler); err != nil {
		return fmt.Errorf("register grpc-gateway handlers: %w", err)
	}

	mergedSwaggerHandler, err := swagger.NewHandler(swaggerMergedOpenAPI)
	if err != nil {
		return fmt.Errorf("init merged swagger-ui handler: %w", err)
	}

	userSwaggerHandler, err := swagger.NewHandler(swaggerUserOpenAPI)
	if err != nil {
		return fmt.Errorf("init user swagger-ui handler: %w", err)
	}

	authSwaggerHandler, err := swagger.NewHandler(swaggerAuthOpenAPI)
	if err != nil {
		return fmt.Errorf("init auth swagger-ui handler: %w", err)
	}

	handler := http.NewServeMux()

	handler.Handle("/", gatewayMux)
	handler.Handle(swaggerBasePath+"/", http.StripPrefix(swaggerBasePath, mergedSwaggerHandler))
	handler.Handle(swaggerBasePath, http.RedirectHandler(swaggerRedirectTarget, http.StatusMovedPermanently))
	handler.Handle(swaggerUserBasePath+"/", http.StripPrefix(swaggerUserBasePath, userSwaggerHandler))
	handler.Handle(swaggerUserBasePath, http.RedirectHandler(swaggerUserRedirectPath, http.StatusMovedPermanently))
	handler.Handle(swaggerAuthBasePath+"/", http.StripPrefix(swaggerAuthBasePath, authSwaggerHandler))
	handler.Handle(swaggerAuthBasePath, http.RedirectHandler(swaggerAuthRedirectPath, http.StatusMovedPermanently))

	c.Gateway = handler

	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
