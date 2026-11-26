package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationUser struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Status         string
	JoinedAt       time.Time
	CreatedAt      time.Time
}
