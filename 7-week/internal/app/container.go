package app

import (
	"context"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/database/postgres"
	gatewayTransport "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/http/gateway"
	grpcTransport "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config/env"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/logger"
)

type Container struct {
	Cfg config.Config

	DB         database.Client
	Tx         database.TxManager
	JWTManager jwt.Manager

	Logger *zap.Logger

	Repositories repository.Repositories
	Services     service.Services

	GRPC          *grpc.Server
	Gateway       http.Handler
	gatewayCancel context.CancelFunc
}

func Build(ctx context.Context, configPath string) (*Container, error) {
	container := &Container{}

	if err := container.initConfig(configPath); err != nil {
		return nil, err
	}
	if err := container.initLogger(); err != nil {
		return nil, err
	}
	if err := container.initDatabase(ctx); err != nil {
		return nil, err
	}
	if err := container.initJWTManager(); err != nil {
		return nil, err
	}

	container.initTxManager()
	container.initRepositories()
	container.initServices()

	if err := container.initGRPC(); err != nil {
		return nil, err
	}
	if err := container.initGateway(ctx); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initConfig(configPath string) error {
	if err := config.Load(configPath); err != nil {
		return fmt.Errorf("load dotenv config: %w", err)
	}

	cfg, err := env.NewConfig()
	if err != nil {
		return fmt.Errorf("load env config: %w", err)
	}

	c.Cfg = cfg

	return nil
}

func (c *Container) initLogger()  error {
	loggerInstance, err := logger.New(logger.Config{
		Level:            c.Cfg.Zap().Level(),
		Encoding:         c.Cfg.Zap().Encoding(),
		OutputPaths:      c.Cfg.Zap().OutputPaths(),
		ErrorOutputPaths: c.Cfg.Zap().ErrorOutputPaths(),
	})
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}

	c.Logger = loggerInstance

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

func (c *Container) initRepositories() {
	c.Repositories = repository.NewRepositories(repository.Deps{DB: c.DB})
}

func (c *Container) initServices() {
	c.Services = service.NewServices(service.Deps{
		Repositories: c.Repositories,
		JWTManager:   c.JWTManager,
	})
}

func (c *Container) initGRPC() error {
	grpcServer, err := grpcTransport.NewServer(grpcTransport.Deps{
		Config: grpcTransport.ServerConfig{
			EnableTLS: c.Cfg.TLS().Enabled(),
			CertFile:  c.Cfg.TLS().CertFile(),
			KeyFile:   c.Cfg.TLS().KeyFile(),
		},
		JWTManager: c.JWTManager,
		Logger:     c.Logger,
		Services:   c.Services,
	})
	if err != nil {
		return fmt.Errorf("init grpc server: %w", err)
	}

	c.GRPC = grpcServer

	return nil
}

func (c *Container) initGateway(ctx context.Context) error {
	handler, cancel, err := gatewayTransport.NewProxy(ctx, gatewayTransport.Config{
		GRPCAddress: c.Cfg.GRPC().Address(),
		TLS:         c.Cfg.TLS(),
	})
	if err != nil {
		return fmt.Errorf("build grpc-gateway proxy: %w", err)
	}

	c.Gateway = handler
	c.gatewayCancel = cancel

	return nil
}

func (c *Container) Close() error {
	if c.gatewayCancel != nil {
		c.gatewayCancel()
	}

	if c.DB != nil {
		return c.DB.Close()
	}

	return nil
}
