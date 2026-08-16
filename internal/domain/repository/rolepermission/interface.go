package rolepermission

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	AssignPermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	RevokePermission(ctx context.Context, roleID, permissionID uuid.UUID) error
	HasPermission(ctx context.Context, roleID, permissionID uuid.UUID) (bool, error)
	ListPermissionsByRole(ctx context.Context, roleID uuid.UUID) ([]*entity.Permission, error)
	ListRolesByPermission(ctx context.Context, permissionID uuid.UUID) ([]*entity.Role, error)
}
