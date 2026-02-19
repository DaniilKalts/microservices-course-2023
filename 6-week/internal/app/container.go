package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/DaniilKalts/microservices-course-2023/6-week/gen/grpc/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/in/database/postgres"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/in/transport/http/swagger"
	grpcInterceptor "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/interceptor"
	userAPI "github.com/DaniilKalts/microservices-course-2023/6-week/internal/adapters/out/transport/grpc/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository"
	userRepository "github.com/DaniilKalts/microservices-course-2023/6-week/internal/repository/user"
	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/service"
	userService "github.com/DaniilKalts/microservices-course-2023/6-week/internal/service/user"
)

const (
	swaggerBasePath       = "/swagger"
	swaggerOpenAPIFile    = "gen/openapi/user/v1/user.swagger.json"
	swaggerRedirectTarget = swaggerBasePath + "/"
)

type Container struct {
	Cfg config.Config

	DB database.Client
	Tx database.TxManager

	UserRepo    repository.UserRepository
	UserSvc     service.UserService
	userHandler userv1.UserV1Server

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

	container.initTxManager()
	container.initUserRepository()
	container.initUserService()
	container.initUserHandler()

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

func (c *Container) initUserService() {
	c.UserSvc = userService.NewService(c.UserRepo)
}

func (c *Container) initUserHandler() {
	c.userHandler = userAPI.NewHandler(c.UserSvc)
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

	c.GRPC = grpcServer

	return nil
}

func (c *Container) initGateway(ctx context.Context) error {
	gatewayMux := runtime.NewServeMux()
	if err := userv1.RegisterUserV1HandlerServer(ctx, gatewayMux, c.userHandler); err != nil {
		return fmt.Errorf("register grpc-gateway handlers: %w", err)
	}

	swaggerHandler, err := swagger.NewHandler(swaggerOpenAPIFile)
	if err != nil {
		return fmt.Errorf("init swagger-ui handler: %w", err)
	}

	handler := http.NewServeMux()

	handler.Handle("/", gatewayMux)
	handler.Handle(swaggerBasePath+"/", http.StripPrefix(swaggerBasePath, swaggerHandler))
	handler.Handle(swaggerBasePath, http.RedirectHandler(swaggerRedirectTarget, http.StatusMovedPermanently))

	c.Gateway = handler

	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
