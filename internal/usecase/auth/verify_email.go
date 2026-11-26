package auth

import (
	"context"
	"time"

	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/emailverificationtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	VerifyEmailParams struct {
		Token string
	}
	VerifyEmailResult struct {
		Message string
	}
)

type verifyEmailUsecase struct {
	verifRepo emailverificationtoken.Repository
	userRepo  user.Repository
	logger    service.Logger
}

func NewVerifyEmailUsecase(
	verifRepo emailverificationtoken.Repository,
	userRepo user.Repository,
	logger service.Logger,
) VerifyEmailUsecase {
	return &verifyEmailUsecase{
		verifRepo: verifRepo,
		userRepo:  userRepo,
		logger:    logger,
	}
}

func (uc *verifyEmailUsecase) Execute(ctx context.Context, params VerifyEmailParams) (*VerifyEmailResult, error) {
	tokenHash := hash.HashSHA256(params.Token)

	token, err := uc.verifRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "verification token is invalid")
	}

	if token.UsedAt != nil {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "verification token has expired")
	}

	if err := uc.verifRepo.MarkUsed(ctx, token.ID); err != nil {
		uc.logger.Error(ctx, "failed to mark verification token as used",
			"action", "VERIFY_EMAIL",
			"user_id", token.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to process verification")
	}

	now := time.Now()
	if err := uc.userRepo.UpdateEmailVerified(ctx, token.UserID, &now); err != nil {
		uc.logger.Error(ctx, "failed to update user email_verified_at",
			"action", "VERIFY_EMAIL",
			"user_id", token.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to verify email")
	}

	if er := uc.userRepo.UpdateStatus(ctx, token.UserID, "active"); er != nil {
		uc.logger.Error(ctx, "failed to update user status to active",
			"action", "VERIFY_EMAIL",
			"user_id", token.UserID,
			"error", er.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to activate user account")
	}

	uc.logger.Info(ctx, "email verified successfully",
		"action", "VERIFY_EMAIL",
		"user_id", token.UserID,
	)

	return &VerifyEmailResult{
		Message: "Email verified successfully",
	}, nil
}
