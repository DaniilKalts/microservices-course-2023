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
	gatewayAddr := a.c.Cfg.Gateway().Address()

	httpSrv := &http.Server{
		Addr:              gatewayAddr,
		Handler:           a.c.Gateway,
		ReadHeaderTimeout: 5 * time.Second,
	}

	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", grpcAddr, err)
	}

	serveErr := make(chan error, 1)

	go func() {
		serveErr <- serveGRPC(a.c.GRPC, grpcListener)
	}()

	go func() {
		serveErr <- serveGateway(httpSrv)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stop)

	select {
	case err := <-serveErr:
		a.shutdown(5*time.Second, httpSrv)
		return err
	case <-ctx.Done():
		a.shutdown(5*time.Second, httpSrv)
		return ctx.Err()
	case <-stop:
		a.shutdown(5*time.Second, httpSrv)
		return nil
	}
}

func serveGRPC(server *grpc.Server, lis net.Listener) error {
	if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("grpc server: %w", err)
	}

	return nil
}

func serveGateway(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("grpc-gateway server: %w", err)
	}

	return nil
}

func (a *App) shutdown(timeout time.Duration, httpSrv *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		a.gracefulStopGRPC(timeout)
	}()

	go func() {
		defer wg.Done()
		_ = httpSrv.Shutdown(ctx)
	}()

	wg.Wait()
}

func (a *App) gracefulStopGRPC(timeout time.Duration) {
	done := make(chan struct{})
	go func() {
		a.c.GRPC.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		a.c.GRPC.Stop()
	}
}
