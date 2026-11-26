package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	ID             uuid.UUID
	OrganizationID *uuid.UUID
	ApplicationID  uuid.UUID
	Name           string
	Slug           string
	Description    *string
	IsSystem       bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
