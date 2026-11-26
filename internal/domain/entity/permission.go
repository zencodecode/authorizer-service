package entity

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	Slug          string
	Resource      string
	Action        string
	Description   *string
	CreatedAt     time.Time
}
