package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/app"
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
	"github.com/DaniilKalts/microservices-course-2023/7-week/pkg/logger"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", ".env", "path to config file")
}

func main() {
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.Load(configPath)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	appLogger, err := logger.New(logger.Config{
		Level:            cfg.Zap.Level,
		Encoding:         cfg.Zap.Encoding,
		OutputPaths:      cfg.Zap.OutputPaths,
		ErrorOutputPaths: cfg.Zap.ErrorOutputPaths,
	})
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}

	a, err := app.New(ctx, cfg, appLogger)
	if err != nil {
		exitWithLoggerError(appLogger, "failed to initialize app", err)
	}
	if err = a.Run(ctx); err != nil {
		exitWithLoggerError(appLogger, "application exited with error", err)
	}

	_ = appLogger.Sync()
}

func exitWithLoggerError(logger *zap.Logger, message string, err error) {
	logger.Error(message, zap.Error(err))
	_ = logger.Sync()
	os.Exit(1)
}
