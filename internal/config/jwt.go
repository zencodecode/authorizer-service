package config

import (
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zencodecode/authorizer-service/pkg/envutil"
	"github.com/zencodecode/authorizer-service/pkg/rsakey"
)

const defaultPrivateKeyPath = "./private.pem"
const defaultPublicKeyPath = "./public.pem"

type JWT struct {
	PrivateKeyPath string
	PublicKeyPath  string
	PrivateKey     *rsa.PrivateKey
	PublicKey      *rsa.PublicKey
	KeyID          string
	TokenExpiry    time.Duration
	RefreshExpiry  time.Duration
}

func LoadJWTConfig() (JWT, error) {
	privateKeyPath := envutil.Get("JWT_PRIVATE_KEY_PATH", defaultPrivateKeyPath)
	publicKeyPath := envutil.Get("JWT_PUBLIC_KEY_PATH", defaultPublicKeyPath)
	tokenExpiry := envutil.GetDuration("JWT_TOKEN_EXPIRY", 15*time.Minute)
	refreshExpiry := envutil.GetDuration("JWT_REFRESH_EXPIRY", 7*24*time.Hour)

	if privateKeyPath == "" || publicKeyPath == "" {
		return JWT{}, fmt.Errorf("JWT_PRIVATE_KEY_PATH and JWT_PUBLIC_KEY_PATH must be set")
	}

	privateKey, err := loadPrivateKeyFromEnvOrFile()
	if err != nil {
		return JWT{}, fmt.Errorf("failed to load private key: %w", err)
	}

	jwt := JWT{
		PrivateKeyPath: privateKeyPath,
		PublicKeyPath:  publicKeyPath,
		PrivateKey:     privateKey,
		PublicKey:      &privateKey.PublicKey,
		KeyID:          generateKID(&privateKey.PublicKey),
		TokenExpiry:    tokenExpiry,
		RefreshExpiry:  refreshExpiry,
	}

	if err := jwt.Validate(); err != nil {
		return JWT{}, fmt.Errorf("invalid jwt config: %w", err)
	}

	return jwt, nil
}

func loadPrivateKeyFromEnvOrFile() (*rsa.PrivateKey, error) {
	if strKey := envutil.Get("JWT_PRIVATE_KEY", ""); strKey != "" {
		return rsakey.ParsePrivateKey([]byte(strKey))
	}

	paths := []string{
		envutil.Get("JWT_PRIVATE_KEY_PATH", defaultPrivateKeyPath),
		defaultPrivateKeyPath,
	}

	for _, path := range paths {
		key, err := rsakey.LoadPrivateKey(path)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) || strings.Contains(err.Error(), "no such file") {
				continue
			}
			return nil, err
		}
		return key, nil
	}

	return nil, fmt.Errorf("private key not found: tried paths %v and JWT_PRIVATE_KEY env", paths)
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
