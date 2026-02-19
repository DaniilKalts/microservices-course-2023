package env

import (
	"errors"
	"os"
	"strconv"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config"
)

const (
	tlsEnabledEnvName  = "TLS_ENABLED"
	tlsCertFileEnvName = "TLS_CERT_FILE"
	tlsKeyFileEnvName  = "TLS_KEY_FILE"
)

type tlsConfig struct {
	enabled  bool
	certFile string
	keyFile  string
}

func NewTLSConfig() (config.TLSConfig, error) {
	enabledRaw := os.Getenv(tlsEnabledEnvName)
	if len(enabledRaw) == 0 {
		return nil, errors.New(tlsEnabledEnvName + " is not set")
	}

	enabled, err := strconv.ParseBool(enabledRaw)
	if err != nil {
		return nil, errors.New(tlsEnabledEnvName + " must be a boolean")
	}

	certFile := os.Getenv(tlsCertFileEnvName)
	keyFile := os.Getenv(tlsKeyFileEnvName)

	if enabled {
		if len(certFile) == 0 {
			return nil, errors.New(tlsCertFileEnvName + " is not set")
		}

		if len(keyFile) == 0 {
			return nil, errors.New(tlsKeyFileEnvName + " is not set")
		}
	}

	return &tlsConfig{
		enabled:  enabled,
		certFile: certFile,
		keyFile:  keyFile,
	}, nil
}

func (cfg *tlsConfig) Enabled() bool {
	return cfg.enabled
}

func (cfg *tlsConfig) CertFile() string {
	return cfg.certFile
}

func (cfg *tlsConfig) KeyFile() string {
	return cfg.keyFile
}
