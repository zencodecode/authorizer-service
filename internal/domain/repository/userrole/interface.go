package userrole

import (
	"context"

	"github.com/zencodecode/authorizer-service/internal/domain/entity"
)

type Repository interface {
	Assign(ctx context.Context, userID, roleID string) error
	Revoke(ctx context.Context, userID, roleID string) error
	Replace(ctx context.Context, userID string, roleIDs []string) error
	GetRolesByUserID(ctx context.Context, userID string) ([]*entity.Role, error)
	GetUsersByRoleID(ctx context.Context, roleID string) ([]*entity.User, error)
	HasRole(ctx context.Context, userID, roleID string) (bool, error)
}
