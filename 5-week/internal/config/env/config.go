package env

import (
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/config"
)

type appConfig struct {
	grpc     config.GRPCConfig
	postgres config.PostgresConfig
	gateway  config.GatewayConfig
}

func NewConfig() (config.Config, error) {
	grpcConfig, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	postgresConfig, err := NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	gatewayConfig, err := NewGatewayConfig()
	if err != nil {
		return nil, err
	}

	return &appConfig{
		grpc:     grpcConfig,
		postgres: postgresConfig,
		gateway:  gatewayConfig,
	}, nil
}

func (cfg *appConfig) GRPC() config.GRPCConfig {
	return cfg.grpc
}

func (cfg *appConfig) Postgres() config.PostgresConfig {
	return cfg.postgres
}

func (cfg *appConfig) Gateway() config.GatewayConfig {
	return cfg.gateway
}
