package rolepermission

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	AssignPermission(ctx context.Context, roleID, permissionID string) error
	RevokePermission(ctx context.Context, roleID, permissionID string) error
	HasPermission(ctx context.Context, roleID, permissionID string) (bool, error)
	ListPermissionsByRole(ctx context.Context, roleID string) ([]*entity.Permission, error)
	ListRolesByPermission(ctx context.Context, permissionID string) ([]*entity.Role, error)
}
