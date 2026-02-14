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
	DB  database.Client
	Tx  database.TxManager

	UserRepo repository.UserRepository
	UserSvc  service.UserService

	GRPC    *grpc.Server
	Gateway http.Handler
}

func Build(ctx context.Context, configPath string) (*Container, error) {
	config.Load(configPath)

	cfg, err := env.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("load env config: %w", err)
	}

	db, err := postgres.New(ctx, cfg.Postgres().DSN())
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}

	tx := transaction.NewTransactionManager(db.DB())

	userRepo := userRepository.NewRepository(db)
	userSvc := userService.NewService(userRepo)
	userHandler := userAPI.NewHandler(userSvc)

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	userv1.RegisterUserV1Server(grpcServer, userHandler)

	gatewayMux := runtime.NewServeMux()
	if err = userv1.RegisterUserV1HandlerServer(ctx, gatewayMux, userHandler); err != nil {
		return nil, fmt.Errorf("register grpc-gateway handlers: %w", err)
	}

	return &Container{
		Cfg:      cfg,
		DB:       db,
		Tx:       tx,
		UserRepo: userRepo,
		UserSvc:  userSvc,
		GRPC:     grpcServer,
		Gateway:  gatewayMux,
	}, nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
