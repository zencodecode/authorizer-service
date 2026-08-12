package entity

import (
	"time"

	"github.com/google/uuid"
)

type OAuthRefreshToken struct {
	ID            uuid.UUID
	AccessTokenID uuid.UUID
	TokenHash     string
	UserAgent     *string
	IPAddress     *string
	ExpiresAt     time.Time
	RevokedAt     *time.Time
	CreatedAt     time.Time
}
