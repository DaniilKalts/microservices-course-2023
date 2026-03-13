package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	grpcTransport "github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/grpc"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/diagnostic"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/adapters/transport/http/gateway"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
)

const grpcGracefulStopHeadroom = 2 * time.Second

type App struct {
	cfg    *config.Config
	logger *zap.Logger

	grpc       *grpc.Server
	gateway    *gateway.Proxy
	diagnostic *http.Server

	shutdownOnce sync.Once
}

func NewApp(ctx context.Context, cfg *config.Config, logger *zap.Logger, c *Container) (_ *App, err error) {
	a := &App{cfg: cfg, logger: logger}
	defer func() {
		if err != nil {
			a.shutdown()
		}
	}()

	if err = a.initGRPC(c); err != nil {
		return nil, err
	}
	if err = a.initGateway(ctx, c); err != nil {
		return nil, err
	}
	a.initDiagnostic(c)

	return a, nil
}

func (a *App) initGRPC(c *Container) error {
	grpcServer, err := grpcTransport.NewServer(grpcTransport.Deps{
		Config: grpcTransport.ServerConfig{
			EnableTLS:      a.cfg.TLS.Enabled,
			CertFile:       a.cfg.TLS.CertFile,
			KeyFile:        a.cfg.TLS.KeyFile,
			RequestTimeout: a.cfg.GRPC.RequestTimeout,
		},
		JWTManager: c.JWTManager,
		Logger:     a.logger.Named("transport.grpc"),
		Services:   c.Services,
		Tracer:     c.Tracer,
		Registry:   c.Registry,
	})
	if err != nil {
		return fmt.Errorf("init grpc server: %w", err)
	}

	a.grpc = grpcServer
	return nil
}

func (a *App) initGateway(ctx context.Context, c *Container) error {
	proxy, err := gateway.NewProxy(ctx, gateway.Config{
		GRPCAddress:    a.cfg.GRPC.Address(),
		GatewayAddress: a.cfg.Gateway.Address(),
		TLS:            a.cfg.TLS,
		Tracer:         c.Tracer,
	})
	if err != nil {
		return fmt.Errorf("init gateway: %w", err)
	}

	a.gateway = proxy
	return nil
}

func (a *App) initDiagnostic(c *Container) {
	var checkers []diagnostic.HealthChecker
	if c.DB != nil {
		checkers = append(checkers, diagnostic.NewPingChecker("postgres", c.DB.DB()))
	}

	a.diagnostic = diagnostic.NewServer(diagnostic.Deps{
		Address:  a.cfg.Diagnostic.Address(),
		Checkers: checkers,
		Registry: c.Registry,
	})
}

func (a *App) Run(ctx context.Context) error {
	defer a.shutdown()

	logger := a.logger.Named("app")

	grpcAddr := a.cfg.GRPC.Address()
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", grpcAddr, err)
	}

	logger.Info("starting servers",
		zap.String("grpc", grpcAddr),
		zap.String("gateway", a.gateway.Addr()),
		zap.String("diagnostic", a.diagnostic.Addr),
	)

	serveErr := make(chan error, 3)

	go func() {
		serveErr <- serveGRPC(a.grpc, grpcListener)
	}()
	go func() {
		serveErr <- a.gateway.Serve()
	}()
	go func() {
		serveErr <- serveHTTP(a.diagnostic)
	}()

	select {
	case err := <-serveErr:
		if err != nil {
			logger.Error("server exited with error", zap.Error(err))
		}
		return err
	case <-ctx.Done():
		logger.Info("shutting down...")
		return ctx.Err()
	}
}

func serveGRPC(server *grpc.Server, lis net.Listener) error {
	if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("grpc server: %w", err)
	}

	return nil
}

func serveHTTP(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server %s: %w", server.Addr, err)
	}

	return nil
}

func (a *App) shutdown() {
	a.shutdownOnce.Do(func() {
		shutdownTimeout := a.cfg.App.ShutdownTimeout

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		var wg sync.WaitGroup

		if a.grpc != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				a.gracefulStopGRPC(shutdownTimeout - grpcGracefulStopHeadroom)
			}()
		}

		if a.gateway != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := a.gateway.Shutdown(ctx); err != nil {
					a.logger.Error("failed to shutdown gateway", zap.Error(err))
				}
			}()
		}

		if a.diagnostic != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := a.diagnostic.Shutdown(ctx); err != nil {
					a.logger.Error("failed to shutdown diagnostic server", zap.Error(err))
				}
			}()
		}

		wg.Wait()
		a.logger.Info("shutdown complete")
	})
}

func (a *App) gracefulStopGRPC(timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		a.grpc.GracefulStop()
		close(done)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
	case <-timer.C:
		a.logger.Warn("grpc graceful stop timeout, forcing", zap.Duration("timeout", timeout))
		a.grpc.Stop()
	}
}
