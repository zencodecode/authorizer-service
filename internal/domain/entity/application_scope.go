package entity

import "github.com/google/uuid"

type ApplicationScope struct {
	ID            uuid.UUID
	ApplicationID uuid.UUID
	Scope         string
	Description   *string
}
