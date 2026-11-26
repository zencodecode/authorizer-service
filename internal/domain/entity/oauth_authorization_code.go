package entity

import (
	"time"

	"github.com/google/uuid"
)

type OAuthAuthorizationCode struct {
	ID                  uuid.UUID
	CodeHash            string
	UserID              uuid.UUID
	OrganizationID      *uuid.UUID
	ApplicationID       uuid.UUID
	RedirectURI         string
	CodeChallenge       string
	CodeChallengeMethod string
	Scopes              []string
	ExpiresAt           time.Time
	UsedAt              *time.Time
	CreatedAt           time.Time
}
