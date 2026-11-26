package entity

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID                   uuid.UUID
	Name                 string
	Slug                 string
	ClientID             string
	ClientSecretHash     string
	RedirectURIs         []string
	AllowedGrantTypes    []string
	RequiresOrganization bool
	Metadata             map[string]any
	IsActive             bool
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
