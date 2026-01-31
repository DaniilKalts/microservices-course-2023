package env

import (
	"errors"
	"net"
	"os"

	"github.com/DaniilKalts/microservices-course-2023/3-week/internal/config"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host string
	port string
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(grpcHostEnvName + " is not set")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New(grpcPortEnvName + " is not set")
	}

	return &grpcConfig{
		host: host,
		port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
