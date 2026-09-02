package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/emailverificationtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
	"github.com/zencodecode/authorizer-service/pkg/randutil"
)

type (
	RegisterParams struct {
		Email    string
		Password string
		Name     string
	}

	RegisterResult struct {
		User *entity.User
	}
)

type registerUsecase struct {
	userRepo    user.Repository
	verifRepo   emailverificationtoken.Repository
	emailSender service.EmailSender
	cfg         *config.Config
	logger      service.Logger
}

func NewRegisterUsecase(
	userRepo user.Repository,
	verifRepo emailverificationtoken.Repository,
	emailSender service.EmailSender,
	cfg *config.Config,
	logger service.Logger,
) RegisterUsecase {
	return &registerUsecase{
		userRepo:    userRepo,
		verifRepo:   verifRepo,
		emailSender: emailSender,
		cfg:         cfg,
		logger:      logger,
	}
}

func (uc *registerUsecase) Execute(ctx context.Context, params RegisterParams) (*RegisterResult, error) {
	existingUser, _ := uc.userRepo.GetByEmail(ctx, params.Email)
	if existingUser != nil {
		uc.logger.Warn(ctx, "registration failed: email already exists",
			"action", "REGISTER",
			"email", params.Email,
		)
		return nil, apperr.NewDirectError(enum.INVALID_REQUEST, "email already registered")
	}

	if len(params.Password) < 8 {
		return nil, apperr.NewDirectError(enum.INVALID_REQUEST, "password must be at least 8 characters")
	}
	hashedPassword, err := hash.Hash(params.Password)
	if err != nil {
		uc.logger.Error(ctx, "failed to hash password",
			"action", "REGISTER",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to process registration")
	}

	now := time.Now()
	user := &entity.User{
		ID:           uuid.Must(uuid.NewV7()),
		Email:        params.Email,
		PasswordHash: hashedPassword,
		Name:         params.Name,
		Status:       "pending_verification",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.logger.Error(ctx, "failed to create user",
			"action", "REGISTER",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to create user")
	}

	token, err := randutil.GenerateRandomString(32)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate verification token",
			"action", "REGISTER",
			"error", err.Error())
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate verification token")
	}

	hashedToken := hash.HashSHA256(token)
	expiresIn := 15 * time.Minute

	// TODO: Generate email verification token and send email
	verifyToken := &entity.EmailVerificationToken{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: now.Add(expiresIn),
		CreatedAt: now,
	}

	if err := uc.verifRepo.Create(ctx, verifyToken); err != nil {
		uc.logger.Error(ctx, "failed to verfy user",
			"action", "REGISTER",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to verfy user")
	}

	verifyURL := fmt.Sprintf("%s/verify-email?token=%s",
		uc.cfg.Auth.OIDC.Issuer, token)

	err = uc.emailSender.Send(ctx, service.SendEmailParams{
		To:      user.Email,
		Subject: "Verify your email",
		Body:    fmt.Sprintf(`<p>Click <a href="%s">here</a> to verify your email.</p>`, verifyURL),
	})

	if err != nil {
		uc.logger.Error(ctx, "failed to send verification email",
			"action", "REGISTER",
			"email", user.Email,
			"error", err.Error(),
		)
	}

	uc.logger.Info(ctx, "user registered successfully",
		"action", "REGISTER",
		"user_id", user.ID,
		"email", user.Email,
	)

	return &RegisterResult{User: user}, nil
}
