package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
	"github.com/zencodecode/authorizer-service/pkg/stringopr"
)

type (
	RegisterParams struct {
		Email    string
		Password string
		FullName string
	}

	RegisterOutput struct {
		User *entity.User
	}
)

type registerUsecase struct {
	userRepo     user.Repository
	roleRepo     role.Repository
	userRoleRepo userrole.Repository
	logger       service.Logger
}

func NewRegisterUsecase(
	userRepo user.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
	logger service.Logger,
) RegisterUsecase {
	return &registerUsecase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		logger:       logger,
	}
}

func (uc *registerUsecase) Execute(ctx context.Context, params RegisterParams) (*RegisterOutput, error) {

	if len(params.Password) < 8 {
		uc.logger.Warn(ctx, "Registration failed: password too short",
			"email", params.Email,
			"context", "REGISTER",
		)
		return nil, errors.New("password must be at least 8 characters")
	}

	existingUser, _ := uc.userRepo.GetByEmail(ctx, params.Email)
	if existingUser != nil {
		uc.logger.Warn(ctx, "Registration failed: email already exists",
			"email", params.Email,
			"context", "REGISTER",
		)
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := hash.Hash(params.Password)
	if err != nil {
		uc.logger.Error(ctx, "Failed to hash password",
			"email", params.Email,
			"error", err.Error(),
			"context", "REGISTER",
		)
		return nil, errors.New("failed to hash password")
	}

	userID := stringopr.GenerateUUID()
	parsedUUID, err := uuid.Parse(userID)
	if err != nil {
		uc.logger.Error(ctx, "failed to parse UUID",
			"user_id", userID,
			"error", err.Error(),
			"context", "REGISTER",
		)
		return nil, errors.New("failed to hash password")
	}

	user := &entity.User{
		ID:           parsedUUID,
		Email:        params.Email,
		PasswordHash: hashedPassword,
		Name:         params.FullName,
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.logger.Error(ctx, "Failed to create user",
			"email", params.Email,
			"error", err.Error(),
			"context", "REGISTER",
		)
		return nil, errors.New("failed to create user")
	}
	return &RegisterOutput{User: user}, nil
}
