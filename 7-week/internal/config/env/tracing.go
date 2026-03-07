package env

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
)

const (
	tracingEnabledEnvName         = "TRACING_ENABLED"
	tracingServiceNameEnvName     = "TRACING_SERVICE_NAME"
	tracingJaegerAgentHostEnvName = "TRACING_JAEGER_AGENT_HOST"
	tracingJaegerAgentPortEnvName = "TRACING_JAEGER_AGENT_PORT"
	tracingSamplerTypeEnvName     = "TRACING_SAMPLER_TYPE"
	tracingSamplerParamEnvName    = "TRACING_SAMPLER_PARAM"

	defaultTracingEnabled         = false
	defaultTracingServiceName     = "api"
	defaultTracingJaegerAgentHost = "localhost"
	defaultTracingJaegerAgentPort = "6831"
	defaultTracingSamplerType     = "const"
	defaultTracingSamplerParam    = 1.0
)

type tracingConfig struct {
	enabled         bool
	serviceName     string
	jaegerAgentHost string
	jaegerAgentPort string
	samplerType     string
	samplerParam    float64
}

func NewTracingConfig() (config.TracingConfig, error) {
	enabled, err := readTracingEnabled()
	if err != nil {
		return nil, err
	}

	serviceName := os.Getenv(tracingServiceNameEnvName)
	if len(serviceName) == 0 {
		serviceName = defaultTracingServiceName
	}

	agentHost := os.Getenv(tracingJaegerAgentHostEnvName)
	if len(agentHost) == 0 {
		agentHost = defaultTracingJaegerAgentHost
	}

	agentPort := os.Getenv(tracingJaegerAgentPortEnvName)
	if len(agentPort) == 0 {
		agentPort = defaultTracingJaegerAgentPort
	}

	samplerType := os.Getenv(tracingSamplerTypeEnvName)
	if len(samplerType) == 0 {
		samplerType = defaultTracingSamplerType
	}

	samplerParam, err := readTracingSamplerParam()
	if err != nil {
		return nil, err
	}

	return &tracingConfig{
		enabled:         enabled,
		serviceName:     serviceName,
		jaegerAgentHost: agentHost,
		jaegerAgentPort: agentPort,
		samplerType:     samplerType,
		samplerParam:    samplerParam,
	}, nil
}

func readTracingEnabled() (bool, error) {
	raw := os.Getenv(tracingEnabledEnvName)
	if len(raw) == 0 {
		return defaultTracingEnabled, nil
	}

	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, errors.New(tracingEnabledEnvName + " must be a boolean")
	}

	return value, nil
}

func readTracingSamplerParam() (float64, error) {
	raw := os.Getenv(tracingSamplerParamEnvName)
	if len(raw) == 0 {
		return defaultTracingSamplerParam, nil
	}

	value, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return 0, fmt.Errorf("%s must be a number: %w", tracingSamplerParamEnvName, err)
	}

	if value < 0 {
		return 0, errors.New(tracingSamplerParamEnvName + " must be non-negative")
	}

	return value, nil
}

func (cfg *tracingConfig) Enabled() bool {
	return cfg.enabled
}

func (cfg *tracingConfig) ServiceName() string {
	return cfg.serviceName
}

func (cfg *tracingConfig) JaegerAgentHostPort() string {
	return net.JoinHostPort(cfg.jaegerAgentHost, cfg.jaegerAgentPort)
}

func (cfg *tracingConfig) SamplerType() string {
	return cfg.samplerType
}

func (cfg *tracingConfig) SamplerParam() float64 {
	return cfg.samplerParam
}
