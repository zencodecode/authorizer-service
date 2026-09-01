package auth

import (
	"context"
	"slices"

	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
)

type (
	UserInfoParams struct {
		UserID uuid.UUID
		Scopes []string
	}

	UserInfoResult struct {
		Sub           string `json:"sub"`
		Name          string `json:"name"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
	}
)

type userInfoUsecase struct {
	userRepo user.Repository
	logger   service.Logger
}

func NewUserInfoUsecase(
	userRepo user.Repository,
	logger service.Logger,
) *userInfoUsecase {
	return &userInfoUsecase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *userInfoUsecase) Execute(ctx context.Context, params UserInfoParams) (*UserInfoResult, error) {
	if !slices.Contains(params.Scopes, "openid") {
		return nil, apperr.NewDirectError(enum.ACCESS_DENIED, "missing required scope: openid")
	}

	user, err := uc.userRepo.GetByID(ctx, params.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "CONSENT",
			"user_id", params.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user")
	}

	if slices.Contains(params.Scopes, "profile") {
		userInfo := &UserInfoResult{
			Sub:           user.ID.String(),
			Name:          user.Name,
			Email:         user.Email,
			EmailVerified: user.EmailVerifiedAt != nil,
		}
		return userInfo, nil
	}

	return &UserInfoResult{
		Sub: user.ID.String(),
	}, nil
}
