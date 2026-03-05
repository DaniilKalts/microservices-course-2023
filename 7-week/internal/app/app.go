package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
)

const (
	gatewayReadHeaderTimeout = 5 * time.Second
	shutdownTimeout          = 5 * time.Second
)

type App struct {
	c *Container
}

func New(ctx context.Context, cfg config.Config, logger *zap.Logger) (*App, error) {
	c, err := Build(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	return &App{c: c}, nil
}

func (a *App) Run(ctx context.Context) error {
	logger := a.Logger().Named("app")

	defer func() {
		if err := a.c.Close(); err != nil {
			logger.Error("failed to close application resources", zap.Error(err))
		}
	}()

	grpcAddr := a.c.Cfg.GRPC().Address()
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", grpcAddr, err)
	}

	logger.Info("grpc listener initialized", zap.String("address", grpcAddr))

	gatewayAddr := a.c.Cfg.Gateway().Address()
	gatewayServer := &http.Server{
		Addr:              gatewayAddr,
		Handler:           a.c.Gateway,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
	}

	prometheusServer := a.c.Prometheus

	logger.Info(
		"starting application servers",
		zap.String("grpc_address", grpcAddr),
		zap.String("gateway_address", gatewayAddr),
		zap.String("prometheus_address", prometheusServer.Addr),
	)

	serveErr := make(chan error, 3)

	go func() {
		logger.Info("grpc server started")
		serveErr <- serveGRPC(a.c.GRPC, grpcListener)
	}()

	go func() {
		logger.Info("http gateway server started")
		serveErr <- serveGateway(gatewayServer, a.c.Cfg.TLS())
	}()

	go func() {
		logger.Info("prometheus server started")
		serveErr <- servePrometheus(prometheusServer)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stop)

	var runErr error
	select {
	case runErr = <-serveErr:
		if runErr != nil {
			logger.Error("application server exited with error", zap.Error(runErr))
		}
	case <-ctx.Done():
		runErr = ctx.Err()
		logger.Warn("application context canceled", zap.Error(runErr))
	case sig := <-stop:
		logger.Info("shutdown signal received", zap.String("signal", sig.String()))
	}

	logger.Info("shutting down application servers")
	a.shutdown(shutdownTimeout, gatewayServer, prometheusServer)
	logger.Info("application servers stopped")

	return runErr
}

func (a *App) Logger() *zap.Logger {
	if a == nil || a.c == nil || a.c.Logger == nil {
		panic("app logger is not initialized")
	}

	return a.c.Logger
}

func serveGRPC(server *grpc.Server, lis net.Listener) error {
	if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("grpc server: %w", err)
	}

	return nil
}

func serveGateway(server *http.Server, tlsCfg config.TLSConfig) (err error) {
	if tlsCfg.Enabled() {
		err = server.ListenAndServeTLS(tlsCfg.CertFile(), tlsCfg.KeyFile())
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("grpc-gateway server: %w", err)
	}

	return nil
}

func servePrometheus(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("prometheus server: %w", err)
	}

	return nil
}

func (a *App) shutdown(timeout time.Duration, gatewayServer *http.Server, prometheusServer *http.Server) {
	logger := a.Logger().Named("shutdown")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		a.gracefulStopGRPC(timeout)
	}()

	go func() {
		defer wg.Done()
		if err := gatewayServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("failed to shutdown gateway server", zap.Error(err))
		}
	}()

	go func() {
		defer wg.Done()
		if err := prometheusServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			logger.Error("failed to shutdown prometheus server", zap.Error(err))
		}
	}()

	wg.Wait()
}

func (a *App) gracefulStopGRPC(timeout time.Duration) {
	logger := a.Logger().Named("shutdown.grpc")

	done := make(chan struct{})
	go func() {
		a.c.GRPC.GracefulStop()
		close(done)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
		logger.Info("grpc server stopped gracefully")
	case <-timer.C:
		logger.Warn("grpc graceful stop timeout reached, forcing stop", zap.Duration("timeout", timeout))
		a.c.GRPC.Stop()
	}
}
