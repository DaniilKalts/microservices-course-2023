package env

import (
	"errors"
	"net"
	"os"
	"strconv"
)

type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	Host string
	Port int
}

func NewGRPCConfig(hostEnvName, portEnvName string) (GRPCConfig, error) {
	host := os.Getenv(hostEnvName)
	if len(host) == 0 {
		return nil, errors.New(hostEnvName + " is not set")
	}

	portStr := os.Getenv(portEnvName)
	if len(portStr) == 0 {
		return nil, errors.New(portEnvName + " is not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, errors.New(portEnvName + " is invalid")
	}

	return &grpcConfig{
		Host: host,
		Port: port,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port))
}