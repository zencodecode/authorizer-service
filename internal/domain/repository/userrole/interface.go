package userrole

import (
	"context"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	AssignRole(ctx context.Context, userRole *entity.UserRole) error
	RevokeRole(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, roleID uuid.UUID) error
	HasRole(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID, roleID uuid.UUID) (bool, error)
	ListRolesByUser(ctx context.Context, userID uuid.UUID, organizationID *uuid.UUID) ([]*entity.Role, error)
	ListUsersByRole(ctx context.Context, roleID uuid.UUID, limit, offset int) ([]*entity.User, error)
	ListPermissionsByUser(ctx context.Context, userID uuid.UUID, applicationID uuid.UUID, organizationID *uuid.UUID) ([]*entity.Permission, error)
}
