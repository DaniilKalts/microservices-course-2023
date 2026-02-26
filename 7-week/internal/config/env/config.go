package env

import (
	"github.com/DaniilKalts/microservices-course-2023/7-week/internal/config"
)

type appConfig struct {
	postgres config.PostgresConfig
	grpc     config.GRPCConfig
	gateway  config.GatewayConfig
	tls      config.TLSConfig
	jwt      config.JWTConfig
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

	tlsConfig, err := NewTLSConfig()
	if err != nil {
		return nil, err
	}

	jwtConfig, err := NewJWTConfig()
	if err != nil {
		return nil, err
	}

	return &appConfig{
		postgres: postgresConfig,
		grpc:     grpcConfig,
		gateway:  gatewayConfig,
		tls:      tlsConfig,
		jwt:      jwtConfig,
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

func (cfg *appConfig) TLS() config.TLSConfig {
	return cfg.tls
}

func (cfg *appConfig) JWT() config.JWTConfig {
	return cfg.jwt
}
