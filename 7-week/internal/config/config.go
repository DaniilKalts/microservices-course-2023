package config

import (
	"fmt"
	"net"
	"time"
)

type Config struct {
	GRPC       GRPCConfig       `envPrefix:"GRPC_"`
	Postgres   PostgresConfig   `envPrefix:"POSTGRES_"`
	Gateway    GatewayConfig    `envPrefix:"GATEWAY_"`
	Prometheus PrometheusConfig `envPrefix:"PROMETHEUS_"`
	TLS        TLSConfig        `envPrefix:"TLS_"`
	JWT        JWTConfig        `envPrefix:"JWT_"`
	Zap        ZapConfig        `envPrefix:"ZAP_"`
	Tracing    TracingConfig    `envPrefix:"TRACING_"`
}

type GRPCConfig struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (cfg *GRPCConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

type PostgresConfig struct {
	Host     string `env:"HOST,required"`
	Port     string `env:"PORT,required"`
	User     string `env:"USER,required"`
	Password string `env:"PASSWORD,required"`
	Name     string `env:"DB,required"`
	SSLMode  string `env:"SSLMODE,required"`
}

func (cfg *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)
}

type GatewayConfig struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (cfg *GatewayConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

type PrometheusConfig struct {
	Host string `env:"HOST,required"`
	Port string `env:"PORT,required"`
}

func (cfg *PrometheusConfig) Address() string {
	return net.JoinHostPort(cfg.Host, cfg.Port)
}

type TLSConfig struct {
	Enabled  bool   `env:"ENABLED" envDefault:"false"`
	CertFile string `env:"CERT_FILE"`
	KeyFile  string `env:"KEY_FILE"`
}

type JWTConfig struct {
	Issuer           string        `env:"ISS,required"`
	Subject          string        `env:"SUB,required"`
	Audience         string        `env:"AUD,required"`
	PrivateKeyFile   string        `env:"PRIVATE_KEY_FILE" envDefault:"build/jwt/rs256_private.pem"`
	PublicKeyFile    string        `env:"PUBLIC_KEY_FILE"  envDefault:"build/jwt/rs256_public.pem"`
	AccessTokenTTL  time.Duration `env:"ACCESS_EXP,required"`
	RefreshTokenTTL time.Duration `env:"REFRESH_EXP,required"`
	NotBeforeOffset time.Duration `env:"NBF,required"`
	IssuedAtOffset  time.Duration `env:"IAT,required"`
}

type ZapConfig struct {
	Level            string `env:"LEVEL"             envDefault:"info"`
	Encoding         string `env:"ENCODING"          envDefault:"json"`
	OutputPaths      string `env:"OUTPUT_PATHS"      envDefault:"stdout"`
	ErrorOutputPaths string `env:"ERROR_OUTPUT_PATHS" envDefault:"stderr"`
}

type TracingConfig struct {
	Enabled         bool    `env:"ENABLED"            envDefault:"false"`
	ServiceName     string  `env:"SERVICE_NAME"       envDefault:"api"`
	JaegerAgentHost string  `env:"JAEGER_AGENT_HOST"  envDefault:"localhost"`
	JaegerAgentPort string  `env:"JAEGER_AGENT_PORT"  envDefault:"6831"`
	SamplerType     string  `env:"SAMPLER_TYPE"       envDefault:"const"`
	SamplerParam    float64 `env:"SAMPLER_PARAM"      envDefault:"1.0"`
}

func (cfg *TracingConfig) JaegerAgentHostPort() string {
	return net.JoinHostPort(cfg.JaegerAgentHost, cfg.JaegerAgentPort)
}
