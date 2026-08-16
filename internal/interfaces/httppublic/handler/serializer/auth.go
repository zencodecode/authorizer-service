package serializer

import "github.com/zencodecode/authorizer-service/internal/infrastructure/auth"

type (
	JWKSResponse struct {
		Keys []auth.JWK `json:"keys"`
	}
)
