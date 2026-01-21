package user

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/service"
	userRepository "github.com/DaniilKalts/microservices-course-2023/3-week/internal/repository/user"
	userService "github.com/DaniilKalts/microservices-course-2023/3-week/internal/service/user"
)

type ServiceProvider interface {
	GetConfig() env.Config

	GetPGPool(ctx context.Context) *pgxpool.Pool

	GetUserRepository(ctx context.Context) repository.UserRepository
	GetUserService(ctx context.Context) service.UserService
}

type serviceProvider struct {
	cfgPath string
	cfg     env.Config

	pgPool *pgxpool.Pool

	userSvc  service.UserService
	userRepo repository.UserRepository
}

func NewServiceProvider(configPath string) ServiceProvider {
	return &serviceProvider{cfgPath: configPath}
}

func (sp *serviceProvider) GetConfig() env.Config {
	if sp.cfg == nil {
		var err error
		if err = config.Load(sp.cfgPath); err != nil {
			log.Fatalf("failed to load config: %v", err)
		}

		sp.cfg, err = env.NewConfig()
		if err != nil {
			log.Fatalf("failed to get env config: %v", err)
		}
	}
	return sp.cfg
}

func (sp *serviceProvider) GetPGPool(ctx context.Context) *pgxpool.Pool {
	if sp.pgPool == nil {
		pool, err := pgxpool.New(ctx, sp.cfg.Postgres().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}
		if err := pool.Ping(ctx); err != nil {
			log.Fatalf("ping error: %v", err)
		}
		sp.pgPool = pool
	}
	return sp.pgPool
}

func (sp *serviceProvider) GetUserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepo == nil {
		sp.userRepo = userRepository.NewRepository(sp.GetPGPool(ctx))
	}
	return sp.userRepo
}

func (sp *serviceProvider) GetUserService(ctx context.Context) service.UserService {
	if sp.userSvc == nil {
		sp.userSvc = userService.NewService(sp.GetUserRepository(ctx))
	}
	return sp.userSvc
}
