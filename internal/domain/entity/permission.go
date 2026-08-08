package entity

import "time"

type Permission struct {
	ID          string
	Code        string
	Description *string
	Version     int
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
