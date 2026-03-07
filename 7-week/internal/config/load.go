package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

func Load(path string) (*Config, error) {
	if err := godotenv.Load(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) validate() error {
	if cfg.TLS.Enabled {
		if cfg.TLS.CertFile == "" {
			return errors.New("TLS_CERT_FILE is required when TLS is enabled")
		}
		if cfg.TLS.KeyFile == "" {
			return errors.New("TLS_KEY_FILE is required when TLS is enabled")
		}
	}

	if cfg.Tracing.Enabled {
		if cfg.Tracing.ServiceName == "" {
			return errors.New("TRACING_SERVICE_NAME is required when tracing is enabled")
		}
	}

	if _, err := os.Stat(cfg.JWT.PrivateKeyFile); err != nil {
		return fmt.Errorf("JWT_PRIVATE_KEY_FILE: %w", err)
	}
	if _, err := os.Stat(cfg.JWT.PublicKeyFile); err != nil {
		return fmt.Errorf("JWT_PUBLIC_KEY_FILE: %w", err)
	}

	return nil
}
