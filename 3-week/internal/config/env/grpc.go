package env

import (
	"errors"
	"net"
	"os"
	"strconv"
)

const (
	grpcHostEnvName = "USER_GRPC_HOST"
	grpcPortEnvName = "USER_GRPC_PORT"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	Host string
	Port int
}

func NewGRPCConfig() (GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New(grpcHostEnvName + " is not set")
	}

	portStr := os.Getenv(grpcPortEnvName)
	if len(portStr) == 0 {
		return nil, errors.New(grpcPortEnvName + " is not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New(grpcPortEnvName + " is invalid")
	}

	return &grpcConfig{
		Host: host,
		Port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}