package env

import (
	"errors"
	"os"
	"time"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config"
	envutils "github.com/DaniilKalts/microservices-course-2023/6-week/pkg/env"
)

const (
	jwtIssuerEnvName           = "JWT_ISS"
	jwtSubjectEnvName          = "JWT_SUB"
	jwtAudienceEnvName         = "JWT_AUD"
	jwtPrivateKeyFileEnvName   = "JWT_PRIVATE_KEY_FILE"
	jwtPublicKeyFileEnvName    = "JWT_PUBLIC_KEY_FILE"
	jwtAccessExpiresAtEnvName  = "JWT_ACCESS_EXP"
	jwtLegacyAccessExpEnvName  = "JWT_EXP"
	jwtRefreshExpiresAtEnvName = "JWT_REFRESH_EXP"
	jwtNotBeforeEnvName        = "JWT_NBF"
	jwtIssuedAtEnvName         = "JWT_IAT"

	defaultJWTPrivateKeyFile = "build/jwt/rs256_private.pem"
	defaultJWTPublicKeyFile  = "build/jwt/rs256_public.pem"
)

type jwtConfig struct {
	issuer           string
	subject          string
	audience         string
	privateKeyFile   string
	publicKeyFile    string
	accessExpiresAt  time.Duration
	refreshExpiresAt time.Duration
	notBefore        time.Duration
	issuedAt         time.Duration
}

func NewJWTConfig() (config.JWTConfig, error) {
	issuer := os.Getenv(jwtIssuerEnvName)
	if len(issuer) == 0 {
		return nil, errors.New(jwtIssuerEnvName + " is not set")
	}

	subject := os.Getenv(jwtSubjectEnvName)
	if len(subject) == 0 {
		return nil, errors.New(jwtSubjectEnvName + " is not set")
	}

	audience := os.Getenv(jwtAudienceEnvName)
	if len(audience) == 0 {
		return nil, errors.New(jwtAudienceEnvName + " is not set")
	}

	privateKeyFile := os.Getenv(jwtPrivateKeyFileEnvName)
	if len(privateKeyFile) == 0 {
		privateKeyFile = defaultJWTPrivateKeyFile
	}

	publicKeyFile := os.Getenv(jwtPublicKeyFileEnvName)
	if len(publicKeyFile) == 0 {
		publicKeyFile = defaultJWTPublicKeyFile
	}

	accessExpiresAt, err := readAccessExpiresAt()
	if err != nil {
		return nil, err
	}

	refreshExpiresAt, err := envutils.ReadDuration(jwtRefreshExpiresAtEnvName)
	if err != nil {
		return nil, err
	}

	notBefore, err := envutils.ReadDuration(jwtNotBeforeEnvName)
	if err != nil {
		return nil, err
	}

	issuedAt, err := envutils.ReadDuration(jwtIssuedAtEnvName)
	if err != nil {
		return nil, err
	}

	return &jwtConfig{
		issuer:           issuer,
		subject:          subject,
		audience:         audience,
		privateKeyFile:   privateKeyFile,
		publicKeyFile:    publicKeyFile,
		accessExpiresAt:  accessExpiresAt,
		refreshExpiresAt: refreshExpiresAt,
		notBefore:        notBefore,
		issuedAt:         issuedAt,
	}, nil
}

func (cfg *jwtConfig) Issuer() string {
	return cfg.issuer
}

func (cfg *jwtConfig) Subject() string {
	return cfg.subject
}

func (cfg *jwtConfig) Audience() string {
	return cfg.audience
}

func (cfg *jwtConfig) PrivateKeyFile() string {
	return cfg.privateKeyFile
}

func (cfg *jwtConfig) PublicKeyFile() string {
	return cfg.publicKeyFile
}

func readAccessExpiresAt() (time.Duration, error) {
	if len(os.Getenv(jwtAccessExpiresAtEnvName)) != 0 {
		return envutils.ReadDuration(jwtAccessExpiresAtEnvName)
	}

	if len(os.Getenv(jwtLegacyAccessExpEnvName)) != 0 {
		return envutils.ReadDuration(jwtLegacyAccessExpEnvName)
	}

	return 0, errors.New(jwtAccessExpiresAtEnvName + " is not set")
}

func (cfg *jwtConfig) AccessExpiresAt() time.Duration {
	return cfg.accessExpiresAt
}

func (cfg *jwtConfig) RefreshExpiresAt() time.Duration {
	return cfg.refreshExpiresAt
}

func (cfg *jwtConfig) NotBefore() time.Duration {
	return cfg.notBefore
}

func (cfg *jwtConfig) IssuedAt() time.Duration {
	return cfg.issuedAt
}
