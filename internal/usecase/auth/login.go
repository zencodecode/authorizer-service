package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	LoginParams struct {
		Email         string
		Password      string
		ApplicationID uuid.UUID
		OrgID         *uuid.UUID
	}

	LoginResult struct {
		User         *entity.User
		AccessToken  string
		RefreshToken string
		ExpiresIn    int
	}
)

type loginUsecase struct {
	userRepo         user.Repository
	userRoleRepo     userrole.Repository
	rolePermRepo     rolepermission.Repository
	refreshTokenRepo oauthrefreshtoken.Repository
	jwtService       service.JWTService
	logger           service.Logger
}

func NewLoginUsecase(
	userRepo user.Repository,
	userRoleRepo userrole.Repository,
	rolePermRepo rolepermission.Repository,
	refreshTokenRepo oauthrefreshtoken.Repository,
	jwtService service.JWTService,
	logger service.Logger,
) LoginUsecase {
	return &loginUsecase{
		userRepo:         userRepo,
		userRoleRepo:     userRoleRepo,
		rolePermRepo:     rolePermRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
		logger:           logger,
	}
}

func (uc *loginUsecase) Execute(ctx context.Context, params LoginParams) (*LoginResult, error) {
	u, err := uc.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		uc.logger.Warn(ctx, "failed to fetch user by email",
			"action", "LOGIN",
			"email", params.Email,
			"error", err.Error(),
		)
		return nil, ErrInvalidCredentials
	}
	if u == nil {
		return nil, ErrInvalidCredentials
	}

	if !hash.CheckHash(u.PasswordHash, params.Password) {
		uc.logger.Warn(ctx, "invalid password",
			"action", "LOGIN",
			"user_id", u.ID,
		)
		return nil, ErrInvalidCredentials
	}

	if u.Status != "active" {
		uc.logger.Warn(ctx, "login attempt on inactive account",
			"action", "LOGIN",
			"user_id", u.ID,
			"status", u.Status,
		)
		return nil, ErrAccountSuspended
	}

	if u.EmailVerifiedAt != nil {
		uc.logger.Warn(ctx, "login attempt on unverified account",
			"action", "LOGIN",
			"user_id", u.ID,
			"status", u.Status,
		)
		return nil, ErrEmailNotVerified
	}

	roles, err := uc.userRoleRepo.ListRolesByUser(ctx, u.ID, params.OrgID)
	if err != nil {
		uc.logger.Error(ctx, "failed to fetch roles",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, errors.New("failed to fetch user roles")
	}

	roleSlugs := make([]string, 0, len(roles))
	permSet := make(map[string]struct{})

	for _, r := range roles {
		if r.ApplicationID != params.ApplicationID {
			continue
		}
		roleSlugs = append(roleSlugs, r.Slug)

		perms, _ := uc.rolePermRepo.ListPermissionsByRole(ctx, r.ID)
		for _, p := range perms {
			permSet[p.Slug] = struct{}{}
		}
	}

	permSlugs := make([]string, 0, len(permSet))
	for k := range permSet {
		permSlugs = append(permSlugs, k)
	}

	now := time.Now()
	expiresIn := 15 * time.Minute

	claims := &entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authorizer-service",
			Subject:   u.ID.String(),
			Audience:  jwt.ClaimStrings{params.ApplicationID.String()},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.Must(uuid.NewV7()).String(),
		},
		Name:        u.Name,
		Email:       u.Email,
		Scopes:      []string{"openid", "profile", "email"},
		Roles:       roleSlugs,
		Permissions: permSlugs,
	}

	if params.OrgID != nil {
		orgIDStr := params.OrgID.String()
		claims.OrgID = &orgIDStr
	}

	accessToken, err := uc.jwtService.GenerateAccessToken(ctx, claims)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate access token",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := uc.jwtService.GenerateRefreshToken()
	if err != nil {
		uc.logger.Error(ctx, "failed to generate refresh token",
			"action", "LOGIN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, errors.New("failed to generate refresh token")
	}

	// TODO: Store refresh token hash via oauthrefreshtoken repository

	return &LoginResult{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int(expiresIn.Seconds()),
	}, nil
}
