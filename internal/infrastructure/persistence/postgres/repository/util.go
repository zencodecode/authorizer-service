package repository

import "github.com/google/uuid"

// parseUUID parses a string to uuid.UUID. Returns uuid.Nil on failure.
func parseUUID(s string) uuid.UUID {
	id, _ := uuid.Parse(s)
	return id
}
