package config

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}

type GatewayConfig interface {
	Address() string
}

type Config interface {
	Postgres() PostgresConfig
	GRPC() GRPCConfig
	Gateway() GatewayConfig
}
