package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/database/postgres"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/tracing"
)

type Container struct {
	DB         database.Client
	JWTManager jwt.Manager
	Tracer     opentracing.Tracer
	Services   service.Services
}

type tracerCloser struct {
	opentracing.Tracer
	io.Closer
}

func NewContainer(ctx context.Context, cfg *config.Config, logger *zap.Logger) (_ *Container, err error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	c := &Container{}
	defer func() {
		if err != nil {
			_ = c.Close()
		}
	}()

	if err = c.initTracer(cfg.Tracing, logger); err != nil {
		return nil, err
	}
	if err = c.initDatabase(ctx, cfg.Postgres, logger); err != nil {
		return nil, err
	}
	if err = c.initJWTManager(cfg.JWT, logger); err != nil {
		return nil, err
	}

	c.initServices(logger)

	return c, nil
}

func (c *Container) Close() error {
	var errs []error

	if c.DB != nil {
		if err := c.DB.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close database: %w", err))
		}
	}

	if closer, ok := c.Tracer.(io.Closer); ok {
		if err := closer.Close(); err != nil {
			errs = append(errs, fmt.Errorf("close tracer: %w", err))
		}
	}

	return errors.Join(errs...)
}

func (c *Container) initServices(logger *zap.Logger) {
	repos := repository.NewRepositories(c.DB, logger.Named("repository"))
	c.Services = service.NewServices(service.Deps{
		Repositories: repos,
		JWTManager:   c.JWTManager,
		Logger:       logger.Named("service"),
	})
}

func (c *Container) initTracer(cfg config.TracingConfig, logger *zap.Logger) error {
	if !cfg.Enabled {
		c.Tracer = opentracing.NoopTracer{}
		return nil
	}

	tracer, closer, err := tracing.NewJaegerTracer(tracing.Config{
		ServiceName:   cfg.ServiceName,
		AgentHostPort: cfg.JaegerAgentHostPort(),
		SamplerType:   cfg.SamplerType,
		SamplerParam:  cfg.SamplerParam,
	})
	if err != nil {
		return fmt.Errorf("init tracer: %w", err)
	}

	c.Tracer = &tracerCloser{Tracer: tracer, Closer: closer}

	logger.Info("tracing initialized",
		zap.String("service_name", cfg.ServiceName),
		zap.String("jaeger_agent", cfg.JaegerAgentHostPort()),
	)

	return nil
}

func (c *Container) initDatabase(ctx context.Context, cfg config.PostgresConfig, logger *zap.Logger) error {
	db, err := postgres.New(ctx, cfg, logger.Named("storage.postgres"))
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

func (c *Container) initJWTManager(cfg config.JWTConfig, logger *zap.Logger) error {
	privateKey, err := jwt.LoadPrivateKey(cfg.PrivateKeyFile)
	if err != nil {
		return fmt.Errorf("load jwt private key: %w", err)
	}

	publicKey, err := jwt.LoadPublicKey(cfg.PublicKeyFile)
	if err != nil {
		return fmt.Errorf("load jwt public key: %w", err)
	}

	jwtManager, err := jwt.NewManager(privateKey, publicKey, jwt.Config{
		Issuer:          cfg.Issuer,
		Subject:         cfg.Subject,
		Audience:        cfg.Audience,
		AccessTokenTTL:  cfg.AccessTokenTTL,
		RefreshTokenTTL: cfg.RefreshTokenTTL,
		NotBeforeOffset: cfg.NotBeforeOffset,
		IssuedAtOffset:  cfg.IssuedAtOffset,
	})
	if err != nil {
		return fmt.Errorf("init jwt manager: %w", err)
	}

	c.JWTManager = jwtManager
	logger.Info("jwt manager initialized")

	return nil
}
