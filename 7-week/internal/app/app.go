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

func New(ctx context.Context, configPath string) (*App, error) {
	c, err := Build(ctx, configPath)
	if err != nil {
		return nil, err
	}

	return &App{c: c}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		_ = a.c.Close()
	}()

	grpcAddr := a.c.Cfg.GRPC().Address()
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", grpcAddr, err)
	}

	gatewayAddr := a.c.Cfg.Gateway().Address()
	gatewayServer := &http.Server{
		Addr:              gatewayAddr,
		Handler:           a.c.Gateway,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
	}

	prometheusServer := a.c.Prometheus

	serveErr := make(chan error, 3)

	go func() {
		serveErr <- serveGRPC(a.c.GRPC, grpcListener)
	}()

	go func() {
		serveErr <- serveGateway(gatewayServer, a.c.Cfg.TLS())
	}()

	go func() {
		serveErr <- servePrometheus(prometheusServer)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stop)

	var runErr error
	select {
	case runErr = <-serveErr:
	case <-ctx.Done():
		runErr = ctx.Err()
	case <-stop:
	}

	a.shutdown(shutdownTimeout, gatewayServer, prometheusServer)

	return runErr
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
		_ = gatewayServer.Shutdown(shutdownCtx)
	}()

	go func() {
		defer wg.Done()
		_ = prometheusServer.Shutdown(shutdownCtx)
	}()

	wg.Wait()
}

func (a *App) gracefulStopGRPC(timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		a.c.GRPC.GracefulStop()
		close(done)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-done:
	case <-timer.C:
		a.c.GRPC.Stop()
	}
}
