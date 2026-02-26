package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	block, err := loadPEMBlock(path, "private key")
	if err != nil {
		return nil, err
	}

	switch block.Type {
	case "RSA PRIVATE KEY":
		key, parseErr := x509.ParsePKCS1PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("parse PKCS#1 private key %q: %w", path, parseErr)
		}

		return key, nil

	case "PRIVATE KEY":
		keyAny, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("parse PKCS#8 private key %q: %w", path, parseErr)
		}

		key, ok := keyAny.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("parse private key %q: key is not RSA", path)
		}

		return key, nil

	default:
		return nil, fmt.Errorf("decode private key PEM %q: unsupported block type %q", path, block.Type)
	}
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	block, err := loadPEMBlock(path, "public key")
	if err != nil {
		return nil, err
	}

	switch block.Type {
	case "PUBLIC KEY":
		keyAny, parseErr := x509.ParsePKIXPublicKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("parse PKIX public key %q: %w", path, parseErr)
		}

		key, ok := keyAny.(*rsa.PublicKey)
		if !ok {
			return nil, fmt.Errorf("parse public key %q: key is not RSA", path)
		}

		return key, nil

	case "RSA PUBLIC KEY":
		key, parseErr := x509.ParsePKCS1PublicKey(block.Bytes)
		if parseErr != nil {
			return nil, fmt.Errorf("parse PKCS#1 public key %q: %w", path, parseErr)
		}

		return key, nil

	default:
		return nil, fmt.Errorf("decode public key PEM %q: unsupported block type %q", path, block.Type)
	}
}

func loadPEMBlock(path, keyKind string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s file %q: %w", keyKind, path, err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, fmt.Errorf("decode %s PEM %q: no PEM block found", keyKind, path)
	}

	return block, nil
}
