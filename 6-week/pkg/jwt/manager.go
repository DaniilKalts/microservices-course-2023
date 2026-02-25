package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

const (
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 7 * 24 * time.Hour
	defaultIssuedAtOffset  = 0 * time.Second
	defaultNotBeforeOffset = 0 * time.Second

	signingMethodAlgorithm = "RS256"

	operationVerify = "verify"
	operationParse  = "parse"
)

type Manager interface {
	GenerateAccessToken(claims Claims) (string, error)
	GenerateRefreshToken(claims Claims) (string, error)
	Verify(tokenString string) (*Claims, error)
	Parse(tokenString string) (*Claims, error)
	AccessTokenTTL() time.Duration
	RefreshTokenTTL() time.Duration
}

type Config struct {
	Issuer          string
	Subject         string
	Audience        string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	IssuedAtOffset  time.Duration
	NotBeforeOffset time.Duration
}

type manager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey

	issuer          string
	subject         string
	audience        string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
	issuedAtOffset  time.Duration
	notBeforeOffset time.Duration
}

func NewManager(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey, cfg Config) (Manager, error) {
	if privateKey == nil {
		return nil, errors.New("private key is required")
	}

	if publicKey == nil {
		return nil, errors.New("public key is required")
	}

	if !isMatchingKeyPair(privateKey, publicKey) {
		return nil, errors.New("private and public keys must belong to the same key pair")
	}

	cfg, err := normalizeConfig(cfg)
	if err != nil {
		return nil, err
	}

	return &manager{
		privateKey:      privateKey,
		publicKey:       publicKey,
		issuer:          cfg.Issuer,
		subject:         cfg.Subject,
		audience:        cfg.Audience,
		accessTokenTTL:  cfg.AccessTokenTTL,
		refreshTokenTTL: cfg.RefreshTokenTTL,
		issuedAtOffset:  cfg.IssuedAtOffset,
		notBeforeOffset: cfg.NotBeforeOffset,
	}, nil
}

func isMatchingKeyPair(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) bool {
	if privateKey.PublicKey.N == nil || publicKey.N == nil {
		return false
	}

	return privateKey.PublicKey.E == publicKey.E && privateKey.PublicKey.N.Cmp(publicKey.N) == 0
}

func normalizeConfig(cfg Config) (Config, error) {
	if cfg.AccessTokenTTL < 0 {
		return Config{}, errors.New("access token ttl must be non-negative")
	}

	if cfg.RefreshTokenTTL < 0 {
		return Config{}, errors.New("refresh token ttl must be non-negative")
	}

	if cfg.NotBeforeOffset < 0 {
		return Config{}, errors.New("not-before offset must be non-negative")
	}

	if cfg.IssuedAtOffset < 0 {
		return Config{}, errors.New("issued-at offset must be non-negative")
	}

	if cfg.AccessTokenTTL == 0 {
		cfg.AccessTokenTTL = defaultAccessTokenTTL
	}

	if cfg.RefreshTokenTTL == 0 {
		cfg.RefreshTokenTTL = defaultRefreshTokenTTL
	}

	if cfg.IssuedAtOffset == 0 {
		cfg.IssuedAtOffset = defaultIssuedAtOffset
	}

	if cfg.NotBeforeOffset == 0 {
		cfg.NotBeforeOffset = defaultNotBeforeOffset
	}

	return cfg, nil
}

func (m *manager) GenerateAccessToken(claims Claims) (string, error) {
	return m.generateToken(claims, m.accessTokenTTL, accessTokenType)
}

func (m *manager) GenerateRefreshToken(claims Claims) (string, error) {
	return m.generateToken(claims, m.refreshTokenTTL, refreshTokenType)
}

func (m *manager) AccessTokenTTL() time.Duration {
	return m.accessTokenTTL
}

func (m *manager) RefreshTokenTTL() time.Duration {
	return m.refreshTokenTTL
}

func (m *manager) generateToken(claims Claims, ttl time.Duration, tokenType string) (string, error) {
	prepared, err := m.prepareClaims(claims, ttl, tokenType)
	if err != nil {
		return "", err
	}

	return m.sign(prepared)
}

func (m *manager) sign(claims Claims) (string, error) {
	token := jwtv5.NewWithClaims(jwtv5.SigningMethodRS256, claims)

	signed, err := token.SignedString(m.privateKey)
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

func (m *manager) Verify(tokenString string) (*Claims, error) {
	claims, err := m.parseClaims(operationVerify, tokenString, m.verifyParserOptions()...)
	if err != nil {
		return nil, err
	}

	if err = validateTokenType(claims.TokenType, ""); err != nil {
		return nil, fmt.Errorf("%s token: %w", operationVerify, err)
	}

	return claims, nil
}

func (m *manager) Parse(tokenString string) (*Claims, error) {
	claims, err := m.parseClaims(
		operationParse,
		tokenString,
		jwtv5.WithValidMethods([]string{signingMethodAlgorithm}),
		jwtv5.WithoutClaimsValidation(),
	)
	if err != nil {
		return nil, err
	}

	if err = m.validateIdentityClaims(operationParse, claims); err != nil {
		return nil, err
	}

	if err = validateTokenType(claims.TokenType, ""); err != nil {
		return nil, fmt.Errorf("%s token: %w", operationParse, err)
	}

	now := time.Now().UTC()

	if claims.NotBefore != nil && now.Before(claims.NotBefore.Time) {
		return nil, fmt.Errorf("%s token: token is not valid yet", operationParse)
	}

	if claims.IssuedAt != nil && now.Before(claims.IssuedAt.Time) {
		return nil, fmt.Errorf("%s token: token issued in the future", operationParse)
	}

	return claims, nil
}

func (m *manager) parseClaims(
	operation string,
	tokenString string,
	parserOptions ...jwtv5.ParserOption,
) (*Claims, error) {
	tokenString = normalizeToken(tokenString)
	if tokenString == "" {
		return nil, fmt.Errorf("%s token: token is empty", operation)
	}

	claims := &Claims{}

	token, err := jwtv5.ParseWithClaims(tokenString, claims, m.keyFunc, parserOptions...)
	if err != nil {
		return nil, fmt.Errorf("%s token: %w", operation, err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("%s token: token is invalid", operation)
	}

	return claims, nil
}

func (m *manager) validateIdentityClaims(operation string, claims *Claims) error {
	if m.issuer != "" && claims.Issuer != m.issuer {
		return fmt.Errorf("%s token: issuer mismatch: got %q", operation, claims.Issuer)
	}

	if m.audience != "" && !containsAudience(claims.Audience, m.audience) {
		return fmt.Errorf("%s token: audience mismatch: got %v", operation, claims.Audience)
	}

	if m.subject != "" && claims.Subject != m.subject {
		return fmt.Errorf("%s token: subject mismatch: got %q", operation, claims.Subject)
	}

	return nil
}

func (m *manager) verifyParserOptions() []jwtv5.ParserOption {
	options := []jwtv5.ParserOption{
		jwtv5.WithValidMethods([]string{signingMethodAlgorithm}),
	}

	if m.issuer != "" {
		options = append(options, jwtv5.WithIssuer(m.issuer))
	}

	if m.audience != "" {
		options = append(options, jwtv5.WithAudience(m.audience))
	}

	if m.subject != "" {
		options = append(options, jwtv5.WithSubject(m.subject))
	}

	return options
}
