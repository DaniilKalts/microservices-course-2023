package env

import (
	"os"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
)

const (
	zapLevelEnvName            = "ZAP_LEVEL"
	zapEncodingEnvName         = "ZAP_ENCODING"
	zapOutputPathsEnvName      = "ZAP_OUTPUT_PATHS"
	zapErrorOutputPathsEnvName = "ZAP_ERROR_OUTPUT_PATHS"

	defaultZapLevel            = "info"
	defaultZapEncoding         = "json"
	defaultZapOutputPaths      = "stdout"
	defaultZapErrorOutputPaths = "stderr"
)

type zapConfig struct {
	level            string
	encoding         string
	outputPaths      string
	errorOutputPaths string
}

func NewZapConfig() (config.ZapConfig, error) {
	level := os.Getenv(zapLevelEnvName)
	if len(level) == 0 {
		level = defaultZapLevel
	}

	encoding := os.Getenv(zapEncodingEnvName)
	if len(encoding) == 0 {
		encoding = defaultZapEncoding
	}

	outputPaths := os.Getenv(zapOutputPathsEnvName)
	if len(outputPaths) == 0 {
		outputPaths = defaultZapOutputPaths
	}

	errorOutputPaths := os.Getenv(zapErrorOutputPathsEnvName)
	if len(errorOutputPaths) == 0 {
		errorOutputPaths = defaultZapErrorOutputPaths
	}

	return &zapConfig{
		level:            level,
		encoding:         encoding,
		outputPaths:      outputPaths,
		errorOutputPaths: errorOutputPaths,
	}, nil
}

func (cfg *zapConfig) Level() string {
	return cfg.level
}

func (cfg *zapConfig) Encoding() string {
	return cfg.encoding
}

func (cfg *zapConfig) OutputPaths() string {
	return cfg.outputPaths
}

func (cfg *zapConfig) ErrorOutputPaths() string {
	return cfg.errorOutputPaths
}
