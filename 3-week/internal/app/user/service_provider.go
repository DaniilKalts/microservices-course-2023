package user

import (
	"context"
	"log"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/clients/database/postgres"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository"
	userRepository "github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository/user"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/service"
	userService "github.com/DaniilKalts/microservices-course-2023/3-week/internal/service/user"
)

type ServiceProvider interface {
	GetConfig() config.Config

	GetDBClient(ctx context.Context) database.Client

	GetUserRepository(ctx context.Context) repository.UserRepository
	GetUserService(ctx context.Context) service.UserService
	Close()
}

type serviceProvider struct {
	cfgPath string
	cfg     config.Config

	dbClient  database.Client
	txManager database.TxManager

	userSvc  service.UserService
	userRepo repository.UserRepository
}

func NewServiceProvider(configPath string) ServiceProvider {
	return &serviceProvider{cfgPath: configPath}
}

func (sp *serviceProvider) GetConfig() config.Config {
	if sp.cfg == nil {
		config.Load(sp.cfgPath)

		var err error
		sp.cfg, err = env.NewConfig()
		if err != nil {
			log.Fatalf("failed to get env config: %v", err)
		}
	}
	return sp.cfg
}

func (sp *serviceProvider) GetDBClient(ctx context.Context) database.Client {
	if sp.dbClient == nil {
		var err error
		sp.dbClient, err = postgres.New(ctx, sp.GetConfig().Postgres().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
	}
	return sp.dbClient
}

func (sp  *serviceProvider) GetTxManager(ctx context.Context) database.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTransactionManager(sp.GetDBClient(ctx).DB())
	}
	return sp.txManager
}

func (sp *serviceProvider) GetUserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepo == nil {
		sp.userRepo = userRepository.NewRepository(sp.GetDBClient(ctx))
	}
	return sp.userRepo
}

func (sp *serviceProvider) GetUserService(ctx context.Context) service.UserService {
	if sp.userSvc == nil {
		sp.userSvc = userService.NewService(sp.GetUserRepository(ctx), sp.GetTxManager(ctx))
	}
	return sp.userSvc
}

func (sp *serviceProvider) Close() {
	if sp.dbClient != nil {
		sp.dbClient.Close()
	}
}
