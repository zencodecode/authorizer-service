package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"math/big"
)

type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

type JWK struct {
	KTY string `json:"kty"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	KID string `json:"kid"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func BuildJWKS(publicKey *rsa.PublicKey, keyID string) *JWKSResponse {
	n := base64.RawURLEncoding.EncodeToString(publicKey.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(publicKey.E)).Bytes())

	return &JWKSResponse{
		Keys: []JWK{{
			KTY: "RSA",
			Use: "sig",
			Alg: "RS256",
			KID: keyID,
			N:   n,
			E:   e,
		}},
	}
}
