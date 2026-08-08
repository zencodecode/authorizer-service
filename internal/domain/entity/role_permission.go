package entity

import "time"

type RolePermission struct {
	RoleID       string
	PermissionID string
	CreatedAt    time.Time
}
