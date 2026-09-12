package seed

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type seederService struct {
	userRepo     user.Repository
	appRepo      application.Repository
	roleRepo     role.Repository
	permRepo     permission.Repository
	rolePermRepo rolepermission.Repository
	userRoleRepo userrole.Repository
	logger       service.Logger
}

func NewSeederService(
	userRepo user.Repository,
	appRepo application.Repository,
	roleRepo role.Repository,
	permRepo permission.Repository,
	rolePermRepo rolepermission.Repository,
	userRoleRepo userrole.Repository,
	logger service.Logger,
) service.SeederService {
	return &seederService{
		userRepo:     userRepo,
		appRepo:      appRepo,
		roleRepo:     roleRepo,
		permRepo:     permRepo,
		rolePermRepo: rolePermRepo,
		userRoleRepo: userRoleRepo,
		logger:       logger,
	}
}

func (uc *seederService) Seed(ctx context.Context, params service.SeederParams) error {
	// ──────── Step 1: Application "AUTHORIZER" ────────
	app, err := uc.seedApplication(ctx, params.AppClientSecret)
	if err != nil {
		return fmt.Errorf("seed application: %w", err)
	}

	// ──────── Step 2: Permissions ────────
	permIDs, err := uc.seedPermissions(ctx, app.ID)
	if err != nil {
		return fmt.Errorf("seed permissions: %w", err)
	}

	// ──────── Step 3: Role "SUPER_ADMIN" ────────
	superAdminRole, err := uc.seedRole(ctx, app.ID)
	if err != nil {
		return fmt.Errorf("seed role: %w", err)
	}

	// ──────── Step 4: RolePermissions ────────
	if err := uc.seedRolePermissions(ctx, superAdminRole.ID, permIDs); err != nil {
		return fmt.Errorf("seed role permissions: %w", err)
	}

	// ──────── Step 5: Admin User ────────
	adminUser, err := uc.seedAdminUser(ctx, params)
	if err != nil {
		return fmt.Errorf("seed admin user: %w", err)
	}

	// ──────── Step 6: UserRole Assignment ────────
	if err := uc.seedUserRole(ctx, adminUser.ID, superAdminRole.ID); err != nil {
		return fmt.Errorf("seed user role: %w", err)
	}

	uc.logger.Info(ctx, "seed completed successfully",
		"user_id", adminUser.ID,
		"email", adminUser.Email,
		"role", "super_admin",
		"app", "authorizer",
		"permissions_count", len(permIDs),
	)

	return nil
}

func (uc *seederService) seedApplication(ctx context.Context, clientSecret string) (*entity.Application, error) {
	existing, err := uc.appRepo.GetBySlug(ctx, "authorizer")
	if err != nil && !errors.Is(err, application.ErrNotFound) {
		return nil, fmt.Errorf("lookup application: %w", err)
	}

	if existing != nil {
		uc.logger.Info(ctx, "application already exists — skipping",
			"slug", existing.Slug, "id", existing.ID)
		return existing, nil
	}

	secretHash, err := hash.Hash(clientSecret)
	if err != nil {
		return nil, fmt.Errorf("hash client secret: %w", err)
	}

	now := time.Now()
	app := &entity.Application{
		ID:                   uuid.Must(uuid.NewV7()),
		Name:                 "Authorizer Service",
		Slug:                 "authorizer",
		ClientID:             "authorizer-service",
		ClientSecretHash:     secretHash,
		RedirectURIs:         []string{},
		AllowedGrantTypes:    []string{"authorization_code", "refresh_token"},
		RequiresOrganization: false,
		Metadata: map[string]any{
			"type":        "internal",
			"description": "Authorizer service internal application",
		},
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := uc.appRepo.Create(ctx, app); err != nil {
		return nil, fmt.Errorf("create application: %w", err)
	}

	uc.logger.Info(ctx, "application created",
		"slug", app.Slug, "id", app.ID, "client_id", app.ClientID)

	return app, nil
}

func (uc *seederService) seedPermissions(ctx context.Context, appID uuid.UUID) ([]uuid.UUID, error) {
	var permIDs []uuid.UUID

	for _, def := range service.DefaultPermissions {
		existing, err := uc.permRepo.GetBySlug(ctx, appID, def.Slug)
		if err != nil && !errors.Is(err, permission.ErrNotFound) {
			return nil, fmt.Errorf("lookup permission %s: %w", def.Slug, err)
		}

		if existing != nil {
			permIDs = append(permIDs, existing.ID)
			continue
		}

		desc := def.Description
		perm := &entity.Permission{
			ID:            uuid.Must(uuid.NewV7()),
			ApplicationID: appID,
			Slug:          def.Slug,
			Resource:      def.Resource,
			Action:        def.Action,
			Description:   &desc,
			CreatedAt:     time.Now(),
		}

		if err := uc.permRepo.Create(ctx, perm); err != nil {
			return nil, fmt.Errorf("create permission %s: %w", def.Slug, err)
		}

		permIDs = append(permIDs, perm.ID)
	}

	uc.logger.Info(ctx, "permissions seeded", "count", len(permIDs))
	return permIDs, nil
}

func (uc *seederService) seedRole(ctx context.Context, appID uuid.UUID) (*entity.Role, error) {
	existing, err := uc.roleRepo.GetBySlug(ctx, nil, appID, "super_admin")
	if err != nil && !errors.Is(err, role.ErrNotFound) {
		return nil, fmt.Errorf("lookup role: %w", err)
	}

	if existing != nil {
		uc.logger.Info(ctx, "role already exists — skipping",
			"slug", existing.Slug, "id", existing.ID)
		return existing, nil
	}

	desc := "Full system access — all permissions on all resources"
	now := time.Now()
	r := &entity.Role{
		ID:             uuid.Must(uuid.NewV7()),
		OrganizationID: nil,
		ApplicationID:  appID,
		Name:           "Super Admin",
		Slug:           "super_admin",
		Description:    &desc,
		IsSystem:       true,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := uc.roleRepo.Create(ctx, r); err != nil {
		return nil, fmt.Errorf("create role: %w", err)
	}

	uc.logger.Info(ctx, "role created", "slug", r.Slug, "id", r.ID)
	return r, nil
}

func (uc *seederService) seedRolePermissions(ctx context.Context, roleID uuid.UUID, permIDs []uuid.UUID) error {
	assigned := 0

	for _, permID := range permIDs {
		has, err := uc.rolePermRepo.HasPermission(ctx, roleID, permID)
		if err != nil && !errors.Is(err, rolepermission.ErrNotFound) {
			return fmt.Errorf("check role permission: %w", err)
		}

		if has {
			continue
		}

		if err := uc.rolePermRepo.AssignPermission(ctx, roleID, permID); err != nil {
			return fmt.Errorf("assign permission %s to role: %w", permID, err)
		}
		assigned++
	}

	uc.logger.Info(ctx, "role permissions seeded",
		"role_id", roleID, "total", len(permIDs), "newly_assigned", assigned)
	return nil
}

func (uc *seederService) seedAdminUser(ctx context.Context, params service.SeederParams) (*entity.User, error) {
	existing, err := uc.userRepo.GetByEmail(ctx, params.AdminEmail)
	if err != nil && !errors.Is(err, user.ErrNotFound) {
		return nil, fmt.Errorf("lookup user: %w", err)
	}

	if existing != nil {
		uc.logger.Info(ctx, "admin user already exists — skipping",
			"user_id", existing.ID, "email", existing.Email)
		return existing, nil
	}

	hashedPassword, err := hash.Hash(params.AdminPassword)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := &entity.User{
		ID:              uuid.Must(uuid.NewV7()),
		Email:           params.AdminEmail,
		PasswordHash:    hashedPassword,
		Name:            params.AdminName,
		Status:          "active",
		EmailVerifiedAt: &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	uc.logger.Info(ctx, "admin user created",
		"user_id", user.ID, "email", user.Email)
	return user, nil
}

func (uc *seederService) seedUserRole(ctx context.Context, userID, roleID uuid.UUID) error {
	hasRole, err := uc.userRoleRepo.HasRole(ctx, userID, nil, roleID)
	if err != nil && !errors.Is(err, userrole.ErrNotFound) {
		return fmt.Errorf("check user role: %w", err)
	}

	if hasRole {
		uc.logger.Info(ctx, "admin already has super_admin role — skipping",
			"user_id", userID, "role_id", roleID)
		return nil
	}

	ur := &entity.UserRole{
		ID:             uuid.Must(uuid.NewV7()),
		OrganizationID: nil,
		UserID:         userID,
		RoleID:         roleID,
		AssignedAt:     time.Now(),
		AssignedBy:     nil,
	}

	if err := uc.userRoleRepo.AssignRole(ctx, ur); err != nil {
		return fmt.Errorf("assign role: %w", err)
	}

	uc.logger.Info(ctx, "super_admin role assigned",
		"user_id", userID, "role_id", roleID)
	return nil
}
