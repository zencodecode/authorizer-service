package userrole

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	AssignRole(ctx context.Context, userRole *entity.UserRole) error
	RevokeRole(ctx context.Context, userID, organizationID, roleID string) error
	HasRole(ctx context.Context, userID, organizationID, roleID string) (bool, error)
	ListRolesByUser(ctx context.Context, userID string, organizationID *string) ([]*entity.Role, error)
	ListUsersByRole(ctx context.Context, roleID string, limit, offset int) ([]*entity.User, error)
	ListPermissionsByUser(ctx context.Context, userID, applicationID string, organizationID *string) ([]*entity.Permission, error)
}
