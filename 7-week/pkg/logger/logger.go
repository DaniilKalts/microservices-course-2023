package logger

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	Level            string
	Encoding         string
	OutputPaths      string
	ErrorOutputPaths string
}

func New(cfg Config) (*zap.Logger, error) {
	level, err := zap.ParseAtomicLevel(cfg.Level)
	if err != nil {
		return nil, fmt.Errorf("parse zap level: %w", err)
	}

	outputPaths := parsePaths(cfg.OutputPaths)
	if len(outputPaths) == 0 {
		outputPaths = []string{"stdout"}
	}

	errorOutputPaths := parsePaths(cfg.ErrorOutputPaths)
	if len(errorOutputPaths) == 0 {
		errorOutputPaths = []string{"stderr"}
	}

	zapConfig := zap.Config{
		Level:            level,
		Development:      false,
		DisableCaller:    false,
		Sampling:         nil,
		Encoding:         cfg.Encoding,
		OutputPaths:      outputPaths,
		ErrorOutputPaths: errorOutputPaths,
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:    "message",
			LevelKey:      "level",
			TimeKey:       "time",
			NameKey:       "logger",
			CallerKey:     "caller",
			StacktraceKey: "stacktrace",
			EncodeLevel:   zapcore.LowercaseLevelEncoder,
			EncodeTime:    zapcore.ISO8601TimeEncoder,
			EncodeCaller:  zapcore.ShortCallerEncoder,
			EncodeName:    zapcore.FullNameEncoder,
		},
	}

	logger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return logger, nil
}

func parsePaths(raw string) []string {
	parts := strings.Split(raw, ",")
	paths := make([]string, 0, len(parts))

	for _, part := range parts {
		path := strings.TrimSpace(part)
		if len(path) == 0 {
			continue
		}

		paths = append(paths, path)
	}

	return paths
}
