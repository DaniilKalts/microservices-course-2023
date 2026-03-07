package tracing

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type Config struct {
	ServiceName   string
	AgentHostPort string
	SamplerType   string
	SamplerParam  float64
}

func NewJaegerTracer(cfg Config) (opentracing.Tracer, io.Closer, error) {
	if cfg.ServiceName == "" {
		return nil, nil, fmt.Errorf("tracing service name is empty")
	}

	if cfg.AgentHostPort == "" {
		return nil, nil, fmt.Errorf("tracing jaeger agent host:port is empty")
	}

	jaegerCfg := jaegercfg.Configuration{
		ServiceName: cfg.ServiceName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  cfg.SamplerType,
			Param: cfg.SamplerParam,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LocalAgentHostPort: cfg.AgentHostPort,
			LogSpans:           false,
		},
	}

	tracer, closer, err := jaegerCfg.NewTracer()
	if err != nil {
		return nil, nil, fmt.Errorf("init jaeger tracer: %w", err)
	}

	return tracer, closer, nil
}
