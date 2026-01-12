package env

import (
	"errors"
	"fmt"
	"os"
)

const (
	postgresHostEnvName     = "POSTGRES_HOST"
	postgresPortEnvName     = "POSTGRES_PORT"
	postgresUserEnvName     = "POSTGRES_USER"
	postgresPasswordEnvName = "POSTGRES_PASSWORD"
	postgresDBEnvName       = "POSTGRES_DB"
	postgresSSLModeEnvName  = "POSTGRES_SSLMODE"
)

type postgresConfig struct {
	dsn string
}

func NewPostgresConfig() (*postgresConfig, error) {
	host := os.Getenv(postgresHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("POSTGRES_HOST is not set")
	}

	port := os.Getenv(postgresPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("POSTGRES_PORT is not set")
	}

	user := os.Getenv(postgresUserEnvName)
	if len(user) == 0 {
		return nil, errors.New("POSTGRES_USER is not set")
	}

	password := os.Getenv(postgresPasswordEnvName)
	if len(password) == 0 {
		return nil, errors.New("POSTGRES_PASSWORD is not set")
	}

	db := os.Getenv(postgresDBEnvName)
	if len(db) == 0 {
		return nil, errors.New("POSTGRES_DB is not set")
	}

	sslMode := os.Getenv(postgresSSLModeEnvName)
	if len(sslMode) == 0 {
		return nil, errors.New("POSTGRES_SSLMODE is not set")
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user,
		password,
		host,
		port,
		db,
		sslMode,
	)
	return &postgresConfig{dsn: dsn}, nil
}