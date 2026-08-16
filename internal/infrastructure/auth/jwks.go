package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

type JWKSService interface {
	GetJWKS(publicKey *rsa.PublicKey, keyID string) (*JWKSResponse, error)
}

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	Kty string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type jwksService struct{}

func NewJWKSService() JWKSService {
	return &jwksService{}
}

func (s *jwksService) GetJWKS(publicKey *rsa.PublicKey, keyID string) (*JWKSResponse, error) {
	if publicKey == nil {
		return nil, ErrNilPublicKey
	}

	nBytes := publicKey.N.Bytes()
	n := base64.RawURLEncoding.EncodeToString(nBytes)

	eBytes := big.NewInt(int64(publicKey.E)).Bytes()
	e := base64.RawURLEncoding.EncodeToString(eBytes)

	jwk := JWK{
		Kty: "RSA",
		Use: "sig",
		Alg: "RS256",
		Kid: keyID,
		N:   n,
		E:   e,
	}

	response := &JWKSResponse{
		Keys: []JWK{jwk},
	}

	return response, nil
}

var ErrNilPublicKey = &jwksError{message: "public key cannot be nil"}

type jwksError struct {
	message string
}

func (e *jwksError) Error() string {
	return e.message
}
