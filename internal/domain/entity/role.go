package entity

import "time"

type Role struct {
	ID          string
	Code        string
	Name        string
	Description *string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
