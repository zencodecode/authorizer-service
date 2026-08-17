package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	RegisterParams struct {
		Email    string
		Password string
		Name     string
	}

	RegisterOutput struct {
		User *entity.User
	}
)

type registerUsecase struct {
	userRepo user.Repository
	logger   service.Logger
}

func NewRegisterUsecase(
	userRepo user.Repository,
	logger service.Logger,
) RegisterUsecase {
	return &registerUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *registerUsecase) Execute(ctx context.Context, params RegisterParams) (*RegisterOutput, error) {
	if len(params.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	existingUser, _ := uc.userRepo.GetByEmail(ctx, params.Email)
	if existingUser != nil {
		uc.logger.Warn(ctx, "registration failed: email already exists",
			"email", params.Email,
		)
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := hash.Hash(params.Password)
	if err != nil {
		uc.logger.Error(ctx, "failed to hash password",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, errors.New("failed to process registration")
	}

	now := time.Now()
	u := &entity.User{
		ID:           uuid.Must(uuid.NewV7()),
		Email:        params.Email,
		PasswordHash: hashedPassword,
		Name:         params.Name,
		Status:       "pending_verification",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = uc.userRepo.Create(ctx, u)
	if err != nil {
		uc.logger.Error(ctx, "failed to create user",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, errors.New("failed to create user")
	}

	// TODO: Generate email verification token and send email

	uc.logger.Info(ctx, "user registered successfully",
		"user_id", u.ID,
		"email", u.Email,
	)

	return &RegisterOutput{User: u}, nil
}
