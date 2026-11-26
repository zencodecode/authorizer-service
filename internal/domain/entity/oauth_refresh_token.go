package entity

import (
	"time"

	"github.com/google/uuid"
)

type OAuthRefreshToken struct {
	ID             uuid.UUID
	TokenHash      string
	UserID         uuid.UUID
	ApplicationID  uuid.UUID
	OrganizationID *uuid.UUID
	Scopes         []string
	UserAgent      *string
	IPAddress      *string
	ExpiresAt      time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
}
