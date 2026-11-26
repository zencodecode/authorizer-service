package entity

import (
	"time"

	"github.com/google/uuid"
)

type Invitation struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	RoleID         *uuid.UUID
	InvitedBy      uuid.UUID
	Token          string
	Status         string
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
