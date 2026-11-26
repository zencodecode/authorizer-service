package config

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"time"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
	"github.com/zencodecode/authorizer-service/pkg/rsakey"
)

type JWT struct {
	PrivateKey    *rsa.PrivateKey
	PublicKey     *rsa.PublicKey
	PrivateKeyPEM string
	PublicKeyPEM  string
	KeyID         string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

func LoadJWTConfig() (JWT, error) {
	privateKeyPEM := envutil.Get("JWT_PRIVATE_KEY", "")
	if privateKeyPEM == "" {
		if path := envutil.Get("JWT_PRIVATE_KEY_PATH", ""); path != "" {
			if data, err := os.ReadFile(path); err == nil {
				privateKeyPEM = string(data)
			}
		} else if data, err := os.ReadFile("./private.pem"); err == nil {
			privateKeyPEM = string(data)
		}
	}

	if privateKeyPEM == "" {
		return JWT{}, fmt.Errorf("JWT_PRIVATE_KEY must be set")
	}

	privateKey, err := rsakey.ParsePrivateKey([]byte(privateKeyPEM))
	if err != nil {
		return JWT{}, fmt.Errorf("failed to load private key: %w", err)
	}

	publicKeyPEM := envutil.Get("JWT_PUBLIC_KEY", "")
	if publicKeyPEM == "" {
		if path := envutil.Get("JWT_PUBLIC_KEY_PATH", ""); path != "" {
			if data, err := os.ReadFile(path); err == nil {
				publicKeyPEM = string(data)
			}
		} else if data, err := os.ReadFile("./public.pem"); err == nil {
			publicKeyPEM = string(data)
		}
	}

	var publicKey *rsa.PublicKey
	if publicKeyPEM != "" {
		publicKey, err = rsakey.ParsePublicKey([]byte(publicKeyPEM))
		if err != nil {
			return JWT{}, fmt.Errorf("failed to parse public key: %w", err)
		}
	} else {
		publicKey = &privateKey.PublicKey
		if pemStr, err := rsakey.EncodePublicKeyToPEM(publicKey); err == nil {
			publicKeyPEM = pemStr
		}
	}

	tokenExpiry := envutil.GetDuration("JWT_TOKEN_EXPIRY", 15*time.Minute)
	refreshExpiry := envutil.GetDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour)

	jwt := JWT{
		PrivateKey:    privateKey,
		PublicKey:     publicKey,
		PrivateKeyPEM: privateKeyPEM,
		PublicKeyPEM:  publicKeyPEM,
		KeyID:         generateKID(publicKey),
		TokenExpiry:   tokenExpiry,
		RefreshExpiry: refreshExpiry,
	}

	if err := jwt.Validate(); err != nil {
		return JWT{}, fmt.Errorf("invalid jwt config: %w", err)
	}

	return jwt, nil
}

func generateKID(pub *rsa.PublicKey) string {
	hash := sha256.Sum256(pub.N.Bytes())
	return base64.RawURLEncoding.EncodeToString(hash[:8])
}

func (j *JWT) Validate() error {
	if j.PrivateKey == nil {
		return fmt.Errorf("private key is required")
	}
	if j.PublicKey == nil {
		return fmt.Errorf("public key is required")
	}
	if j.TokenExpiry <= 0 {
		return fmt.Errorf("token expiry must be greater than zero")
	}
	if j.RefreshExpiry <= 0 {
		return fmt.Errorf("refresh expiry must be greater than zero")
	}
	return nil
}
