package token

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/zencodecode/authorizer-service/internal/definition/enum"
	"github.com/zencodecode/authorizer-service/internal/domain/apperr"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthaccesstoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthrefreshtoken"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/organization"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	ExchangeParams struct {
		GrantType    string
		Code         string
		RedirectURI  string
		ClientID     string
		ClientSecret string
		CodeVerifier string
	}
	ExchangeResult struct {
		User         *entity.User
		AccessToken  string
		RefreshToken string
		Scope        string
		IDToken      string
	}
)

type exchangeUsecase struct {
	userRepo      user.Repository
	appRepo       application.Repository
	oauthCodeRepo oauthauthorizationcode.Repository
	orgRepo       organization.Repository
	roleRepo      role.Repository
	userRoleRepo  userrole.Repository
	permRepo      permission.Repository
	rolePermRepo  rolepermission.Repository
	oauthAccRepo  oauthaccesstoken.Repository
	oauthRefRepo  oauthrefreshtoken.Repository
	jwtSvc        service.JWTService
	logger        service.Logger
}

func NewExchangeUsecase(
	userRepo user.Repository,
	appRepo application.Repository,
	oauthCodeRepo oauthauthorizationcode.Repository,
	orgRepo organization.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
	permRepo permission.Repository,
	rolePermRepo rolepermission.Repository,
	oauthAccRepo oauthaccesstoken.Repository,
	oauthRefRepo oauthrefreshtoken.Repository,
	jwtSvc service.JWTService,
	logger service.Logger,
) ExchangeUsecase {
	return &exchangeUsecase{
		userRepo:      userRepo,
		appRepo:       appRepo,
		oauthCodeRepo: oauthCodeRepo,
		orgRepo:       orgRepo,
		roleRepo:      roleRepo,
		userRoleRepo:  userRoleRepo,
		permRepo:      permRepo,
		rolePermRepo:  rolePermRepo,
		oauthAccRepo:  oauthAccRepo,
		oauthRefRepo:  oauthRefRepo,
		jwtSvc:        jwtSvc,
		logger:        logger,
	}
}

func (uc *exchangeUsecase) Execute(ctx context.Context, params ExchangeParams) (*ExchangeResult, error) {

	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "TOKEN",
			"client_id", &params.ClientID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query application")
	}

	if !hash.CheckHash(app.ClientSecretHash, params.ClientSecret) {
		return nil, apperr.NewDirectError(enum.INVALID_CLIENT, "client secret is invalid")
	}

	codeHash := hash.HashSHA256(params.Code)

	oauthcode, err := uc.oauthCodeRepo.GetByCodeHash(ctx, codeHash)
	if err != nil {
		uc.logger.Error(ctx, "failed to query oauth authorization code",
			"action", "TOKEN",
			"code_hash", codeHash,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query oauth authorization code")
	}

	if !oauthcode.ExpiresAt.After(time.Now()) {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "authorization code has expired")
	}

	if oauthcode.UsedAt != nil {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "authorization code has already been used")
	}

	if oauthcode.RedirectURI != params.RedirectURI {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "redirect_uri does not match")
	}

	computed := hash.HashSHA256(params.CodeVerifier)
	if computed != oauthcode.CodeChallenge {
		return nil, apperr.NewDirectError(enum.INVALID_GRANT, "PKCE verification failed")
	}

	if err := uc.oauthCodeRepo.MarkUsed(ctx, oauthcode.ID); err != nil {
		uc.logger.Error(ctx, "failed to update oauth_authorization_code used_at",
			"action", "TOKEN",
			"code_id", oauthcode.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to update oauth_authorization_code used_at")
	}

	u, err := uc.userRepo.GetByID(ctx, oauthcode.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "TOKEN",
			"user_id", oauthcode.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user")
	}

	org, err := uc.orgRepo.GetByID(ctx, *oauthcode.OrganizationID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query organization",
			"action", "TOKEN",
			"organization_id", oauthcode.OrganizationID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query organization")
	}

	roles, err := uc.userRoleRepo.ListRolesByUser(ctx, u.ID, &org.ID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query roles",
			"action", "TOKEN",
			"user_id", u.ID,
			"error", err.Error(),
		)
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

	now := time.Now()
	expiresIn := 15 * time.Minute

	claims := &entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authorizer-service",
			Subject:   u.ID.String(),
			Audience:  jwt.ClaimStrings{app.ClientID},
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

	accessToken, err := uc.jwtSvc.GenerateAccessToken(ctx, claims)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate access token",
			"action", "TOKEN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate access token")
	}

	access := &entity.OAuthAccessToken{
		ID:             uuid.Must(uuid.NewV7()),
		TokenHash:      hash.HashSHA256(accessToken),
		UserID:         u.ID,
		OrganizationID: &org.ID,
		ApplicationID:  app.ID,
		Scopes:         oauthcode.Scopes,
		ExpiresAt:      claims.ExpiresAt.Time,
	}

	if err := uc.oauthAccRepo.Create(ctx, access); err != nil {
		uc.logger.Error(ctx, "failed to persist access token",
			"action", "TOKEN",
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to persist access token")
	}

	refreshToken, err := uc.jwtSvc.GenerateRefreshToken()
	if err != nil {
		uc.logger.Error(ctx, "failed to generate refresh token",
			"action", "TOKEN",
			"user_id", u.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate refresh token")
	}

	refresh := &entity.OAuthRefreshToken{
		ID:            uuid.Must(uuid.NewV7()),
		AccessTokenID: access.ID,
		TokenHash:     hash.HashSHA256(refreshToken),
		ExpiresAt:     now.AddDate(0, 0, 30),
	}

	if err := uc.oauthRefRepo.Create(ctx, refresh); err != nil {
		uc.logger.Error(ctx, "failed to persist refresh token",
			"action", "TOKEN",
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to persist refresh token")
	}

	token := &ExchangeResult{
		User:         u,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scope:        "",
		IDToken:      refresh.AccessTokenID.String(),
	}

	return token, nil
}
