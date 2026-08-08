package rolepermission

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Grant(ctx context.Context, roleID, permissionID string) error
	Revoke(ctx context.Context, roleID, permissionID string) error
	Replace(ctx context.Context, roleID string, permissionIDs []string) error
	GetPermissionsByRoleID(ctx context.Context, roleID string) ([]*entity.Permission, error)
}
