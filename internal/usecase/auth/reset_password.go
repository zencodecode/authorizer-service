package auth

import (
	"context"
	"time"

	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/passwordresettoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	ResetPasswordParams struct {
		Token       string
		NewPassword string
	}
	ResetPasswordResult struct {
		Message string
	}
)

type resetPasswordUsecase struct {
	userRepo     user.Repository
	resetRepo    passwordresettoken.Repository
	oauthRefRepo oauthrefreshtoken.Repository
	logger       service.Logger
}

func NewResetPasswordUsecase(
	userRepo user.Repository,
	resetRepo passwordresettoken.Repository,
	oauthRefRepo oauthrefreshtoken.Repository,
	logger service.Logger,
) ResetPasswordUsecase {
	return &resetPasswordUsecase{
		userRepo:     userRepo,
		resetRepo:    resetRepo,
		oauthRefRepo: oauthRefRepo,
		logger:       logger,
	}
}

func (uc *resetPasswordUsecase) Execute(ctx context.Context, params ResetPasswordParams) (*ResetPasswordResult, error) {
	if len(params.NewPassword) < 8 {
		return nil, apperr.NewDirectError(enum.INVALID_REQUEST, "password must be at least 8 characters")
	}

	tokenHash := hash.HashSHA256(params.Token)
	token, err := uc.resetRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "reset token is invalid")
	}

	if token.UsedAt != nil {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "reset token has already been used")
	}

	if !token.ExpiresAt.After(time.Now()) {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "reset token has expired")
	}

	if err := uc.resetRepo.MarkUsed(ctx, token.ID); err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to process reset")
	}

	hashedPassword, err := hash.Hash(params.NewPassword)
	if err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to hash password")
	}

	if err := uc.userRepo.UpdatePassword(ctx, token.UserID, hashedPassword); err != nil {
		uc.logger.Error(ctx, "failed to update password",
			"action", "RESET_PASSWORD",
			"user_id", token.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to update password")
	}

	if err := uc.oauthRefRepo.RevokeAllByUser(ctx, token.UserID); err != nil {
		uc.logger.Error(ctx, "failed to revoke refresh token",
			"action", "REVOKE",
			"error", err.Error())
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to revoke refresh token")
	}

	uc.logger.Info(ctx, "password reset successfully",
		"action", "RESET_PASSWORD",
		"user_id", token.UserID,
	)

	return &ResetPasswordResult{
		Message: "Password reset successfully",
	}, nil
}
