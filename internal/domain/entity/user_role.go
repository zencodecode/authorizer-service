package entity

import (
	"time"

	"github.com/google/uuid"
)

type UserRole struct {
	ID             uuid.UUID
	OrganizationID *uuid.UUID
	UserID         uuid.UUID
	RoleID         uuid.UUID
	AssignedAt     time.Time
	AssignedBy     *uuid.UUID
}
