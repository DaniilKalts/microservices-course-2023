package config

type PostgresConfig interface {
	DSN() string
}

type GRPCConfig interface {
	Address() string
}

type Config interface {
	GRPC() GRPCConfig
	Postgres() PostgresConfig
}
