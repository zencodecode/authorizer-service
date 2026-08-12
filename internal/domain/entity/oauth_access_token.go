package entity

import (
	"time"

	"github.com/google/uuid"
)

type OAuthAccessToken struct {
	ID             uuid.UUID
	TokenHash      string
	UserID         uuid.UUID
	OrganizationID *uuid.UUID
	ApplicationID  uuid.UUID
	Scopes         []string
	ExpiresAt      time.Time
	RevokedAt      *time.Time
	CreatedAt      time.Time
}
