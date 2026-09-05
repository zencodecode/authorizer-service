package token

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	RefreshParams struct {
		GrantType    string
		RefreshToken string
		ClientID     string
		ClientSecret string
	}
	RefreshResult struct {
		AccessToken  string
		RefreshToken string
		Scope        string
		IDToken      string
		ExpiresIn    int
	}
)

type refreshUsecase struct {
	userRepo     user.Repository
	appRepo      application.Repository
	orgRepo      organization.Repository
	userRoleRepo userrole.Repository
	rolePermRepo rolepermission.Repository
	permRepo     permission.Repository
	oauthRefRepo oauthrefreshtoken.Repository
	jwtSvc       service.JWTService
	issuerURL    string
	logger       service.Logger
}

func NewRefreshUsecase(
	userRepo user.Repository,
	appRepo application.Repository,
	orgRepo organization.Repository,
	userRoleRepo userrole.Repository,
	rolePermRepo rolepermission.Repository,
	permRepo permission.Repository,
	oauthRefRepo oauthrefreshtoken.Repository,
	jwtSvc service.JWTService,
	issuerURL string,
	logger service.Logger,
) RefreshUsecase {
	return &refreshUsecase{
		userRepo:     userRepo,
		appRepo:      appRepo,
		orgRepo:      orgRepo,
		userRoleRepo: userRoleRepo,
		rolePermRepo: rolePermRepo,
		permRepo:     permRepo,
		oauthRefRepo: oauthRefRepo,
		jwtSvc:       jwtSvc,
		issuerURL:    issuerURL,
		logger:       logger,
	}
}

func (uc *refreshUsecase) Execute(ctx context.Context, params RefreshParams) (*RefreshResult, error) {
	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		return nil, apperr.NewDirectError(enum.INVALID_CLIENT, "invalid client_id")
	}
	if !hash.CheckHash(app.ClientSecretHash, params.ClientSecret) {
		return nil, apperr.NewDirectError(enum.INVALID_CLIENT, "invalid client credentials")
	}

	tokenHash := hash.HashSHA256(params.RefreshToken)
	oldRefresh, err := uc.oauthRefRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		uc.logger.Warn(ctx, "refresh token not found or expired",
			"action", "REFRESH", "client_id", params.ClientID, "error", err.Error())
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "refresh token is invalid or expired")
	}

	if oldRefresh.RevokedAt != nil {
		uc.logger.Warn(ctx, "attempt to use revoked refresh token",
			"action", "REFRESH", "token_id", oldRefresh.ID)
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "refresh token has been revoked")
	}

	if !oldRefresh.ExpiresAt.After(time.Now()) {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "refresh token has expired")
	}

	u, err := uc.userRepo.GetByID(ctx, oldRefresh.UserID)
	if err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user")
	}
	if u.Status != "active" {
		return nil, apperr.NewDirectError(enum.ACCESS_DENIED, "user account is not active")
	}

	var org *entity.Organization
	var orgID *uuid.UUID
	if oldRefresh.OrganizationID != nil {
		org, err = uc.orgRepo.GetByID(ctx, *oldRefresh.OrganizationID)
		if err != nil {
			return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query organization")
		}
		orgID = &org.ID
	}

	roles, err := uc.userRoleRepo.ListRolesByUser(ctx, u.ID, orgID)
	if err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user roles")
	}

	roleSlugs := make([]string, 0, len(roles))
	permSet := make(map[string]struct{})
	for _, r := range roles {
		if r.ApplicationID != app.ID {
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

	if err := uc.oauthRefRepo.Revoke(ctx, oldRefresh.ID); err != nil {
		uc.logger.Error(ctx, "failed to revoke old refresh token",
			"action", "REFRESH", "error", err.Error())
	}

	now := time.Now()
	expiresIn := 15 * time.Minute
	claims := &entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    uc.issuerURL,
			Subject:   u.ID.String(),
			Audience:  jwt.ClaimStrings{app.ClientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.Must(uuid.NewV7()).String(),
		},
		Scopes:      oldRefresh.Scopes,
		Roles:       roleSlugs,
		Permissions: permSlugs,
	}

	if slices.Contains(oldRefresh.Scopes, "profile") {
		claims.Name = u.Name
	}

	if slices.Contains(oldRefresh.Scopes, "email") {
		claims.Email = u.Email
	}

	if org != nil {
		orgIDStr := org.ID.String()
		claims.OrgID = &orgIDStr
		claims.OrgSlug = &org.Slug
	}

	accessToken, err := uc.jwtSvc.GenerateAccessToken(ctx, claims)
	if err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate access token")
	}

	newRefreshTokenStr, err := uc.jwtSvc.GenerateRefreshToken()
	if err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate refresh token")
	}

	newRefresh := &entity.OAuthRefreshToken{
		ID:             uuid.Must(uuid.NewV7()),
		TokenHash:      hash.HashSHA256(newRefreshTokenStr),
		UserID:         u.ID,
		ApplicationID:  app.ID,
		OrganizationID: orgID,
		Scopes:         oldRefresh.Scopes,
		ExpiresAt:      now.AddDate(0, 0, 30),
		CreatedAt:      now,
	}

	if err := uc.oauthRefRepo.Create(ctx, newRefresh); err != nil {
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to persist refresh token")
	}

	var idToken string
	if slices.Contains(oldRefresh.Scopes, "openid") {
		idClaims := &entity.IDTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    uc.issuerURL,
				Subject:   u.ID.String(),
				Audience:  jwt.ClaimStrings{app.ClientID},
				ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        uuid.Must(uuid.NewV7()).String(),
			},
			AuthTime: now.Unix(),
		}
		if slices.Contains(oldRefresh.Scopes, "profile") {
			idClaims.Name = u.Name
		}
		if slices.Contains(oldRefresh.Scopes, "email") {
			idClaims.Email = u.Email
			idClaims.EmailVerified = u.EmailVerifiedAt != nil
		}
		idToken, _ = uc.jwtSvc.GenerateIDToken(ctx, idClaims)
	}

	return &RefreshResult{
		AccessToken:  accessToken,
		RefreshToken: newRefreshTokenStr,
		Scope:        strings.Join(oldRefresh.Scopes, " "),
		IDToken:      idToken,
		ExpiresIn:    int(expiresIn.Seconds()),
	}, nil
}
