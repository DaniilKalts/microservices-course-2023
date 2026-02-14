package env

import (
	"errors"
	"net"
	"os"

	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/config"
)

const (
	gatewayHostEnvName = "GATEWAY_HOST"
	gatewayPortEnvName  = "GATEWAY_PORT"
)

type gatewayConfig struct {
	host string
	port string
}

func NewGatewayConfig() (config.GRPCConfig, error) {
	host := os.Getenv(gatewayHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(gatewayHostEnvName + " is not set")
	}

	port := os.Getenv(gatewayPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(gatewayPortEnvName + " is not set")
	}

	return &gatewayConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *gatewayConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
