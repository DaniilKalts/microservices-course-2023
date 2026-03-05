package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/database/postgres"
	gatewayTransport "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/http/gateway"
	prometheusTransport "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/in/transport/http/prometheus"
	grpcTransport "github.com/DaniilKalts/microservices-course-2023/7-week/internal/adapters/out/transport/grpc"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/clients/database/transaction"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/jwt"
)

type Container struct {
	Cfg config.Config

	DB         database.Client
	Tx         database.TxManager
	JWTManager jwt.Manager

	Logger *zap.Logger

	Repositories repository.Repositories
	Services     service.Services

	GRPC       *grpc.Server
	Gateway    http.Handler
	Prometheus *http.Server
}

func Build(ctx context.Context, cfg config.Config, logger *zap.Logger) (*Container, error) {
	if cfg == nil {
		return nil, errors.New("app config is nil")
	}

	if logger == nil {
		return nil, errors.New("app logger is nil")
	}

	container := &Container{
		Cfg:    cfg,
		Logger: logger,
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
	if err := container.initPrometheus(); err != nil {
		return nil, err
	}

	return container, nil
}

func (c *Container) initDatabase(ctx context.Context) error {
	logger := c.Logger.Named("di.container.database")
	logger.Info("connecting to postgres")

	db, err := postgres.New(ctx, c.Cfg.Postgres().DSN(), c.Logger.Named("storage.postgres"))
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	if err = db.DB().Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	c.DB = db
	logger.Info("postgres connected")

	return nil
}

func (c *Container) initTxManager() {
	logger := c.Logger.Named("di.container.tx")
	logger.Info("initializing transaction manager")

	c.Tx = transaction.NewTransactionManager(c.DB.DB())

	logger.Info("transaction manager initialized")
}

func (c *Container) initJWTManager() error {
	logger := c.Logger.Named("di.container.jwt")

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
	logger.Info("jwt manager initialized")

	return nil
}

func (c *Container) initRepositories() {
	logger := c.Logger.Named("di.container.repository")

	c.Repositories = repository.NewRepositories(repository.Deps{
		DB:     c.DB,
		Logger: c.Logger.Named("repository"),
	})
	logger.Info("repositories initialized")
}

func (c *Container) initServices() {
	logger := c.Logger.Named("di.container.service")

	c.Services = service.NewServices(service.Deps{
		Repositories: c.Repositories,
		JWTManager:   c.JWTManager,
		Logger:       c.Logger.Named("service"),
	})
	logger.Info("services initialized")
}

func (c *Container) initGRPC() error {
	logger := c.Logger.Named("di.container.transport.grpc")

	grpcServer, err := grpcTransport.NewServer(grpcTransport.Deps{
		Config: grpcTransport.ServerConfig{
			EnableTLS: c.Cfg.TLS().Enabled(),
			CertFile:  c.Cfg.TLS().CertFile(),
			KeyFile:   c.Cfg.TLS().KeyFile(),
		},
		JWTManager: c.JWTManager,
		Logger:     c.Logger.Named("transport.grpc"),
		Services:   c.Services,
	})
	if err != nil {
		return fmt.Errorf("init grpc server: %w", err)
	}

	c.GRPC = grpcServer
	logger.Info("grpc transport initialized")

	return nil
}

func (c *Container) initGateway(ctx context.Context) error {
	logger := c.Logger.Named("di.container.transport.http.gateway")

	handler, err := gatewayTransport.NewProxy(ctx, gatewayTransport.Config{
		GRPCAddress: c.Cfg.GRPC().Address(),
		TLS:         c.Cfg.TLS(),
	})
	if err != nil {
		return fmt.Errorf("build grpc-gateway proxy: %w", err)
	}

	c.Gateway = handler
	logger.Info("http gateway initialized")

	return nil
}

func (c *Container) initPrometheus() error {
	logger := c.Logger.Named("di.container.transport.http.prometheus")

	server, err := prometheusTransport.NewServer(prometheusTransport.Config{
		Address: c.Cfg.Prometheus().Address(),
	})
	if err != nil {
		return fmt.Errorf("init prometheus server: %w", err)
	}

	c.Prometheus = server
	logger.Info("prometheus server initialized")

	return nil
}

func (c *Container) Close() error {
	if c.Logger != nil {
		c.Logger.Named("di.container").Info("closing application resources")
	}

	var closeErr error

	if gatewayCloser, ok := c.Gateway.(interface{ Close() error }); ok {
		if err := gatewayCloser.Close(); err != nil {
			closeErr = err
			if c.Logger != nil {
				c.Logger.Named("di.container.transport.http.gateway").Error("failed to close gateway grpc connection", zap.Error(err))
			}
		}
	}

	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			closeErr = err
			if c.Logger != nil {
				c.Logger.Named("di.container.database").Error("failed to close database client", zap.Error(err))
			}
		}
	}

	return closeErr
}
