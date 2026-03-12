package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/app"
	"github.com/DaniilKalts/microservices-course-2023/8-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/8-week/pkg/logger"
)

func main() {
	configPath := flag.String("config-path", ".env", "path to config file")

	flag.Parse()

	if err := run(*configPath); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "application error: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	appLogger, err := logger.New(logger.Config{
		Level:            cfg.Zap.Level,
		Encoding:         cfg.Zap.Encoding,
		OutputPaths:      cfg.Zap.OutputPaths,
		ErrorOutputPaths: cfg.Zap.ErrorOutputPaths,
	})
	if err != nil {
		return fmt.Errorf("initialize logger: %w", err)
	}
	defer func() { _ = appLogger.Sync() }()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	container, err := app.NewContainer(ctx, cfg, appLogger)
	if err != nil {
		return fmt.Errorf("initialize container: %w", err)
	}
	defer func() {
		if closeErr := container.Close(); closeErr != nil {
			appLogger.Error("container close error", zap.Error(closeErr))
		}
	}()

	a, err := app.NewApp(ctx, cfg, appLogger, container)
	if err != nil {
		return fmt.Errorf("initialize app: %w", err)
	}

	return a.Run(ctx)
}
