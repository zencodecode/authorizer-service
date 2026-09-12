package service

import "context"

type permissionDef struct {
	Slug        string
	Resource    string
	Action      string
	Description string
}

// defaultPermissions adalah daftar semua permissions yang di-seed untuk Application "AUTHORIZER".
var DefaultPermissions = []permissionDef{
	// User
	{Slug: "user.create", Resource: "user", Action: "create", Description: "Create new users"},
	{Slug: "user.read", Resource: "user", Action: "read", Description: "View user details"},
	{Slug: "user.update", Resource: "user", Action: "update", Description: "Update user information"},
	{Slug: "user.delete", Resource: "user", Action: "delete", Description: "Delete users"},
	{Slug: "user.assign_roles", Resource: "user", Action: "assign_roles", Description: "Assign roles to users"},
	// Role
	{Slug: "role.create", Resource: "role", Action: "create", Description: "Create new roles"},
	{Slug: "role.read", Resource: "role", Action: "read", Description: "View role details"},
	{Slug: "role.update", Resource: "role", Action: "update", Description: "Update role information"},
	{Slug: "role.delete", Resource: "role", Action: "delete", Description: "Delete roles"},
	// Permission
	{Slug: "permission.create", Resource: "permission", Action: "create", Description: "Create new permissions"},
	{Slug: "permission.read", Resource: "permission", Action: "read", Description: "View permission details"},
	{Slug: "permission.sync", Resource: "permission", Action: "sync", Description: "Sync permissions from external service"},
	// Application
	{Slug: "application.create", Resource: "application", Action: "create", Description: "Register new applications"},
	{Slug: "application.read", Resource: "application", Action: "read", Description: "View application details"},
	{Slug: "application.update", Resource: "application", Action: "update", Description: "Update application information"},
	// Organization
	{Slug: "organization.create", Resource: "organization", Action: "create", Description: "Create new organizations"},
	{Slug: "organization.read", Resource: "organization", Action: "read", Description: "View organization details"},
	{Slug: "organization.update", Resource: "organization", Action: "update", Description: "Update organization information"},
}

type SeederParams struct {
	AdminEmail      string
	AdminPassword   string
	AdminName       string
	AppClientSecret string
}

type SeederService interface {
	Seed(ctx context.Context, params SeederParams) error
}
