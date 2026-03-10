package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
)

func (a *App) Run(ctx context.Context) error {
	logger := a.logger.Named("app")

	ctx, stop := signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	grpcAddr := a.cfg.GRPC.Address()
	grpcListener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return fmt.Errorf("listen grpc %s: %w", grpcAddr, err)
	}

	logger.Info("starting servers",
		zap.String("grpc", grpcAddr),
		zap.String("gateway", a.gatewayServer.Addr),
		zap.String("diagnostic", a.diagnostic.Addr),
	)

	serveErr := make(chan error, 3)

	go func() { serveErr <- serveGRPC(a.grpc, grpcListener) }()
	go func() { serveErr <- serveGateway(a.gatewayServer, a.cfg.TLS) }()
	go func() { serveErr <- serveHTTP(a.diagnostic) }()

	select {
	case err = <-serveErr:
		if err != nil {
			logger.Error("server exited with error", zap.Error(err))
		}
		return err
	case <-ctx.Done():
		logger.Info("shutting down...")
		return nil
	}
}

func serveGRPC(server *grpc.Server, lis net.Listener) error {
	if err := server.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return fmt.Errorf("grpc server: %w", err)
	}

	return nil
}

func serveGateway(server *http.Server, tlsCfg config.TLSConfig) error {
	var err error
	if tlsCfg.Enabled {
		err = server.ListenAndServeTLS(tlsCfg.CertFile, tlsCfg.KeyFile)
	} else {
		err = server.ListenAndServe()
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("gateway server: %w", err)
	}

	return nil
}

func serveHTTP(server *http.Server) error {
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("http server %s: %w", server.Addr, err)
	}

	return nil
}
