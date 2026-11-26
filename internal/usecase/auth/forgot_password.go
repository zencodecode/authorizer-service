package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/passwordresettoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
	"github.com/zencodecode/authorizer-service/pkg/randutil"
)

type (
	ForgotPasswordParams struct {
		Email string
	}
	ForgotPasswordResult struct {
		Message string
	}
)

type forgotPasswordUsecase struct {
	userRepo    user.Repository
	resetRepo   passwordresettoken.Repository
	emailSender service.EmailSender
	cfg         *config.Config
	logger      service.Logger
}

func NewForgotPasswordUsecase(
	userRepo user.Repository,
	resetRepo passwordresettoken.Repository,
	emailSender service.EmailSender,
	cfg *config.Config,
	logger service.Logger,
) ForgotPasswordUsecase {
	return &forgotPasswordUsecase{
		userRepo:    userRepo,
		resetRepo:   resetRepo,
		emailSender: emailSender,
		cfg:         cfg,
		logger:      logger,
	}
}

func (uc *forgotPasswordUsecase) Execute(ctx context.Context, params ForgotPasswordParams) (*ForgotPasswordResult, error) {
	result := &ForgotPasswordResult{
		Message: "if the email exists, a password reset link has been sent",
	}

	user, err := uc.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		return result, nil
	}

	token, err := randutil.GenerateRandomString(32)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate reset token",
			"action", "FORGOT_PASSWORD",
			"error", err.Error(),
		)
		return result, nil
	}

	hashedToken := hash.HashSHA256(token)
	expiresIn := 15 * time.Minute
	now := time.Now()

	resetToken := &entity.PasswordResetToken{
		ID:        uuid.Must(uuid.NewV7()),
		UserID:    user.ID,
		TokenHash: hashedToken,
		ExpiresAt: now.Add(expiresIn),
		CreatedAt: now,
	}

	if err := uc.resetRepo.Create(ctx, resetToken); err != nil {
		uc.logger.Error(ctx, "failed to persist reset token",
			"action", "FORGOT_PASSWORD",
			"error", err.Error(),
		)
		return result, nil
	}

	resetURL := fmt.Sprintf("%s/password/reset?token=%s",
		uc.cfg.Auth.OIDC.Issuer, token)
	body := fmt.Sprintf(`<p>Click <a href="%s">here</a> to reset your password. This link expires in 1 hour.</p>`, resetURL)

	if err := uc.emailSender.Send(ctx, user.Email, "Reset your password", body); err != nil {
		uc.logger.Error(ctx, "failed to send reset email",
			"action", "FORGOT_PASSWORD",
			"email", user.Email,
			"error", err.Error(),
		)
	}

	return result, nil
}
