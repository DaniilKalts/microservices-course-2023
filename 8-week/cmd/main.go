package main

import (
	"context"
	"flag"
	"fmt"
	"os"

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

	ctx := context.Background()

	a, err := app.New(ctx, cfg, appLogger)
	if err != nil {
		return fmt.Errorf("initialize app: %w", err)
	}
	defer func() { _ = a.Close() }()

	return a.Run(ctx)
}
