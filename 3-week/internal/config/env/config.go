package env

type Config interface {
	GRPC() GRPCConfig
	Postgres() PostgresConfig
}

type AppConfig struct {
	grpc     GRPCConfig
	postgres PostgresConfig
}

func NewConfig() (*AppConfig, error) {
	grpcConfig, err := NewGRPCConfig()
	if err != nil {
		return nil, err
	}

	postgresConfig, err := NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	return &AppConfig{
		grpc:     grpcConfig,
		postgres: postgresConfig,
	}, nil
}

func (cfg *AppConfig) GRPC() GRPCConfig {
	return cfg.grpc
}

func (cfg *AppConfig) Postgres() PostgresConfig {
	return cfg.postgres
}
