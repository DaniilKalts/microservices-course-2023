package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	userv1 "github.com/DaniilKalts/microservices-course-2023/5-week/gen/go/user/v1"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/adapters/in/database/postgres"
	grpcInterceptor "github.com/DaniilKalts/microservices-course-2023/5-week/internal/adapters/out/transport/grpc/interceptor"
	userAPI "github.com/DaniilKalts/microservices-course-2023/5-week/internal/adapters/out/transport/grpc/user"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/repository"
	userRepository "github.com/DaniilKalts/microservices-course-2023/5-week/internal/repository/user"
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/service"
	userService "github.com/DaniilKalts/microservices-course-2023/5-week/internal/service/user"
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
	container.initGRPC()

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

func (c *Container) initGRPC() {
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcInterceptor.ValidationInterceptor(),
		),
	)

	reflection.Register(grpcServer)
	userv1.RegisterUserV1Server(grpcServer, c.userHandler)

	c.GRPC = grpcServer
}

func (c *Container) initGateway(ctx context.Context) error {
	gatewayMux := runtime.NewServeMux()
	if err := userv1.RegisterUserV1HandlerServer(ctx, gatewayMux, c.userHandler); err != nil {
		return fmt.Errorf("register grpc-gateway handlers: %w", err)
	}

	c.Gateway = gatewayMux

	return nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
