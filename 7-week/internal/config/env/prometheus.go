package env

import (
	"errors"
	"net"
	"os"

	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
)

const (
	prometheusHostEnvName = "PROMETHEUS_HOST"
	prometheusPortEnvName = "PROMETHEUS_PORT"
)

type prometheusConfig struct {
	host string
	port string
}

func NewPrometheusConfig() (config.PrometheusConfig, error) {
	host := os.Getenv(prometheusHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(prometheusHostEnvName + " is not set")
	}

	port := os.Getenv(prometheusPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(prometheusPortEnvName + " is not set")
	}

	return &prometheusConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *prometheusConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
