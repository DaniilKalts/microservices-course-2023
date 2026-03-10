package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/database/postgres"
	grpcTransport "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/diagnostic"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/gateway"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/clients/database"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/repository"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/service"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/jwt"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/tracing"
)

const (
	gatewayReadHeaderTimeout = 5 * time.Second
	shutdownTimeout          = 5 * time.Second
	grpcGracefulStopTimeout  = 3 * time.Second
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger

	db         database.Client
	jwtManager jwt.Manager
	tracer     opentracing.Tracer
	tracerStop io.Closer

	grpc          *grpc.Server
	gateway       *gateway.Proxy
	gatewayServer *http.Server
	diagnostic    *http.Server
}

func New(ctx context.Context, cfg *config.Config, logger *zap.Logger) (*App, error) {
	if cfg == nil {
		return nil, errors.New("config is nil")
	}
	if logger == nil {
		return nil, errors.New("logger is nil")
	}

	a := &App{cfg: cfg, logger: logger}

	if err := a.init(ctx); err != nil {
		_ = a.Close()
		return nil, err
	}

	return a, nil
}

func (a *App) init(ctx context.Context) error {
	if err := a.initTracer(); err != nil {
		return err
	}
	if err := a.initDatabase(ctx); err != nil {
		return err
	}
	if err := a.initJWTManager(); err != nil {
		return err
	}

	repos := repository.NewRepositories(a.db, a.logger.Named("repository"))
	services := service.NewServices(service.Deps{
		Repositories: repos,
		JWTManager:   a.jwtManager,
		Logger:       a.logger.Named("service"),
	})

	if err := a.initGRPC(services); err != nil {
		return err
	}
	if err := a.initGateway(ctx); err != nil {
		return err
	}
	a.initDiagnostic()

	return nil
}

func (a *App) initTracer() error {
	tracingCfg := a.cfg.Tracing

	if !tracingCfg.Enabled {
		a.tracer = opentracing.GlobalTracer()
		return nil
	}

	tracer, closer, err := tracing.NewJaegerTracer(tracing.Config{
		ServiceName:   tracingCfg.ServiceName,
		AgentHostPort: tracingCfg.JaegerAgentHostPort(),
		SamplerType:   tracingCfg.SamplerType,
		SamplerParam:  tracingCfg.SamplerParam,
	})
	if err != nil {
		return fmt.Errorf("init tracer: %w", err)
	}

	opentracing.SetGlobalTracer(tracer)

	a.tracer = tracer
	a.tracerStop = closer

	a.logger.Info("tracing initialized",
		zap.String("service_name", tracingCfg.ServiceName),
		zap.String("jaeger_agent", tracingCfg.JaegerAgentHostPort()),
	)

	return nil
}

func (a *App) initDatabase(ctx context.Context) error {
	db, err := postgres.New(ctx, a.cfg.Postgres, a.logger.Named("storage.postgres"))
	if err != nil {
		return fmt.Errorf("connect postgres: %w", err)
	}

	if err = db.DB().Ping(ctx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}

	a.db = db
	a.logger.Info("postgres connected")

	return nil
}

func (a *App) initJWTManager() error {
	privateKey, err := jwt.LoadPrivateKey(a.cfg.JWT.PrivateKeyFile)
	if err != nil {
		return fmt.Errorf("load jwt private key: %w", err)
	}

	publicKey, err := jwt.LoadPublicKey(a.cfg.JWT.PublicKeyFile)
	if err != nil {
		return fmt.Errorf("load jwt public key: %w", err)
	}

	jwtManager, err := jwt.NewManager(privateKey, publicKey, jwt.Config{
		Issuer:          a.cfg.JWT.Issuer,
		Subject:         a.cfg.JWT.Subject,
		Audience:        a.cfg.JWT.Audience,
		AccessTokenTTL:  a.cfg.JWT.AccessTokenTTL,
		RefreshTokenTTL: a.cfg.JWT.RefreshTokenTTL,
		NotBeforeOffset: a.cfg.JWT.NotBeforeOffset,
		IssuedAtOffset:  a.cfg.JWT.IssuedAtOffset,
	})
	if err != nil {
		return fmt.Errorf("init jwt manager: %w", err)
	}

	a.jwtManager = jwtManager
	a.logger.Info("jwt manager initialized")

	return nil
}

func (a *App) initGRPC(services service.Services) error {
	grpcServer, err := grpcTransport.NewServer(grpcTransport.Deps{
		Config: grpcTransport.ServerConfig{
			EnableTLS:      a.cfg.TLS.Enabled,
			CertFile:       a.cfg.TLS.CertFile,
			KeyFile:        a.cfg.TLS.KeyFile,
			RequestTimeout: a.cfg.GRPC.Timeout,
		},
		JWTManager: a.jwtManager,
		Logger:     a.logger.Named("transport.grpc"),
		Services:   services,
		Tracer:     a.tracer,
	})
	if err != nil {
		return fmt.Errorf("init grpc server: %w", err)
	}

	a.grpc = grpcServer
	a.logger.Info("grpc server initialized")

	return nil
}

func (a *App) initGateway(ctx context.Context) error {
	proxy, err := gateway.NewProxy(ctx, gateway.Config{
		GRPCAddress: a.cfg.GRPC.Address(),
		TLS:         a.cfg.TLS,
	})
	if err != nil {
		return fmt.Errorf("build grpc-gateway proxy: %w", err)
	}

	a.gateway = proxy
	a.gatewayServer = &http.Server{
		Addr:              a.cfg.Gateway.Address(),
		Handler:           proxy,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
	}
	a.logger.Info("gateway initialized")

	return nil
}

func (a *App) initDiagnostic() {
	a.diagnostic = diagnostic.NewServer(diagnostic.Deps{
		Address: a.cfg.Prometheus.Address(),
		DB:      a.db,
	})
	a.logger.Info("diagnostic server initialized")
}

func (a *App) Close() error {
	a.stopServers()

	var errs []error

	if a.gateway != nil {
		if err := a.gateway.Close(); err != nil {
			a.logger.Error("failed to close gateway", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if a.db != nil {
		if err := a.db.Close(); err != nil {
			a.logger.Error("failed to close database", zap.Error(err))
			errs = append(errs, err)
		}
	}

	if a.tracerStop != nil {
		if err := a.tracerStop.Close(); err != nil {
			a.logger.Error("failed to close tracer", zap.Error(err))
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (a *App) stopServers() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var wg sync.WaitGroup

	if a.grpc != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.gracefulStopGRPC()
		}()
	}

	if a.gatewayServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := a.gatewayServer.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
				a.logger.Error("failed to shutdown gateway server", zap.Error(err))
			}
		}()
	}

	if a.diagnostic != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := a.diagnostic.Shutdown(ctx); err != nil && !errors.Is(err, context.Canceled) {
				a.logger.Error("failed to shutdown diagnostic server", zap.Error(err))
			}
		}()
	}

	wg.Wait()
}

func (a *App) gracefulStopGRPC() {
	done := make(chan struct{})
	go func() {
		a.grpc.GracefulStop()
		close(done)
	}()

	timer := time.NewTimer(grpcGracefulStopTimeout)
	defer timer.Stop()

	select {
	case <-done:
	case <-timer.C:
		a.logger.Warn("grpc graceful stop timeout, forcing", zap.Duration("timeout", grpcGracefulStopTimeout))
		a.grpc.Stop()
	}
}
