package app

import (
	"context"
	"fmt"

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

	GRPC *grpc.Server
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

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	userv1.RegisterUserV1Server(grpcServer, userAPI.NewHandler(userSvc))

	return &Container{
		Cfg:      cfg,
		DB:       db,
		Tx:       tx,
		UserRepo: userRepo,
		UserSvc:  userSvc,
		GRPC:     grpcServer,
	}, nil
}

func (c *Container) Close() error {
	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
