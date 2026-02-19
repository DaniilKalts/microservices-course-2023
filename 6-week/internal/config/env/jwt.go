package env

import (
	"errors"
	"os"

	"github.com/DaniilKalts/microservices-course-2023/6-week/internal/config"
)

const (
	jwtIssuerEnvName    = "JWT_ISS"
	jwtSubjectEnvName   = "JWT_SUB"
	jwtAudienceEnvName  = "JWT_AUD"
	jwtExpiresAtEnvName = "JWT_EXP"
	jwtNotBeforeEnvName = "JWT_NBF"
	jwtIssuedAtEnvName  = "JWT_IAT"
)

type jwtConfig struct {
	issuer    string
	subject   string
	audience  string
	expiresAt string
	notBefore string
	issuedAt  string
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

	expiresAt := os.Getenv(jwtExpiresAtEnvName)
	if len(expiresAt) == 0 {
		return nil, errors.New(jwtExpiresAtEnvName + " is not set")
	}

	notBefore := os.Getenv(jwtNotBeforeEnvName)
	if len(notBefore) == 0 {
		return nil, errors.New(jwtNotBeforeEnvName + " is not set")
	}

	issuedAt := os.Getenv(jwtIssuedAtEnvName)
	if len(issuedAt) == 0 {
		return nil, errors.New(jwtIssuedAtEnvName + " is not set")
	}

	return &jwtConfig{
		issuer:    issuer,
		subject:   subject,
		audience:  audience,
		expiresAt: expiresAt,
		notBefore: notBefore,
		issuedAt:  issuedAt,
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

func (cfg *jwtConfig) ExpiresAt() string {
	return cfg.expiresAt
}

func (cfg *jwtConfig) NotBefore() string {
	return cfg.notBefore
}

func (cfg *jwtConfig) IssuedAt() string {
	return cfg.issuedAt
}
