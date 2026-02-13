package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	c   *Container
	lis net.Listener
}

func New(ctx context.Context, configPath string) (*App, error) {
	c, err := Build(ctx, configPath)
	if err != nil {
		return nil, err
	}

	return &App{c: c}, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() { _ = a.c.Close() }()

	addr := a.c.Cfg.GRPC().Address()
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("listen %s: %w", addr, err)
	}
	a.lis = lis

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- a.c.GRPC.Serve(lis)
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	defer signal.Stop(stop)

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		a.gracefulStopWithTimeout(5 * time.Second)
		return ctx.Err()
	case <-stop:
		a.gracefulStopWithTimeout(5 * time.Second)
		return nil
	}
}

func (a *App) gracefulStopWithTimeout(d time.Duration) {
	done := make(chan struct{})
	go func() {
		a.c.GRPC.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(d):
		a.c.GRPC.Stop()
	}
}
