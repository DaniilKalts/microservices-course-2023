package env

import (
	"github.com/DaniilKalts/microservices-course-2023/5-week/internal/config"
)

type appConfig struct {
	grpc     config.GRPCConfig
	postgres config.PostgresConfig
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

	return &appConfig{
		grpc:     grpcConfig,
		postgres: postgresConfig,
	}, nil
}

func (cfg *appConfig) GRPC() config.GRPCConfig {
	return cfg.grpc
}

func (cfg *appConfig) Postgres() config.PostgresConfig {
	return cfg.postgres
}
