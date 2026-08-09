package entity

import (
	"time"

	"github.com/google/uuid"
)

type Application struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description *string
	Metadata    map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
