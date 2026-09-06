package discovery

import (
	"crypto/rsa"
)

type Handler struct {
	publicKey *rsa.PublicKey
	keyID     string
	issuer    string
}

func New(
	publicKey *rsa.PublicKey,
	keyID string,
	issuer string,
) *Handler {
	return &Handler{
		publicKey: publicKey,
		keyID:     keyID,
		issuer:    issuer,
	}
}
