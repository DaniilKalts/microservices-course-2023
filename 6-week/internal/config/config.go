package config

import "time"

type GRPCConfig interface {
	Address() string
}

type PostgresConfig interface {
	DSN() string
}

type GatewayConfig interface {
	Address() string
}

type TLSConfig interface {
	Enabled() bool
	CertFile() string
	KeyFile() string
}

type JWTConfig interface {
	Issuer() string
	Subject() string
	Audience() string
	AccessExpiresAt() time.Duration
	RefreshExpiresAt() time.Duration
	NotBefore() time.Duration
	IssuedAt() time.Duration
}

type Config interface {
	Postgres() PostgresConfig
	GRPC() GRPCConfig
	Gateway() GatewayConfig
	TLS() TLSConfig
	JWT() JWTConfig
}
