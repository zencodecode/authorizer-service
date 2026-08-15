package serializer

import "github.com/zencodecode/authorizer-service/internal/infrastructure/driver/auth"

type (
	JWKSResponse struct {
		Keys []auth.JWK `json:"keys"`
	}
)
