package env

import (
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/config"
)

type appConfig struct {
	postgres config.PostgresConfig
	grpc     config.GRPCConfig
	gateway  config.GatewayConfig
}

func NewConfig() (config.Config, error) {
	postgresConfig, err := NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	grpcConfig, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	gatewayConfig, err := NewGatewayConfig()
	if err != nil {
		return nil, err
	}

	return &appConfig{
		postgres: postgresConfig,
		grpc:     grpcConfig,
		gateway:  gatewayConfig,
	}, nil
}

func (cfg *appConfig) Postgres() config.PostgresConfig {
	return cfg.postgres
}

func (cfg *appConfig) GRPC() config.GRPCConfig {
	return cfg.grpc
}

func (cfg *appConfig) Gateway() config.GatewayConfig {
	return cfg.gateway
}
