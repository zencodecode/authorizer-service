package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationApplication struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	ApplicationID  uuid.UUID
	IsActive       bool
	ActivatedAt    time.Time
}
