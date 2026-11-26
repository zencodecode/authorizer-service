package entity

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID             uuid.UUID
	OrganizationID *uuid.UUID
	ApplicationID  *uuid.UUID
	ActorUserID    *uuid.UUID
	Action         string
	ResourceType   *string
	ResourceID     *uuid.UUID
	Metadata       map[string]any
	CreatedAt      time.Time
}
