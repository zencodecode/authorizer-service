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
		ExpiresIn    int
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
	oauthRefRepo  oauthrefreshtoken.Repository
	jwtSvc        service.JWTService
	logger        service.Logger
	issuerURL     string
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
	oauthRefRepo oauthrefreshtoken.Repository,
	jwtSvc service.JWTService,
	logger service.Logger,
	issuerURL string,
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
		oauthRefRepo:  oauthRefRepo,
		jwtSvc:        jwtSvc,
		logger:        logger,
		issuerURL:     issuerURL,
	}
}

func (uc *exchangeUsecase) Execute(ctx context.Context, params ExchangeParams) (*ExchangeResult, error) {

	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "EXCHANGE",
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
			"action", "EXCHANGE",
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
			"action", "EXCHANGE",
			"code_id", oauthcode.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to update oauth_authorization_code used_at")
	}

	user, err := uc.userRepo.GetByID(ctx, oauthcode.UserID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query user",
			"action", "EXCHANGE",
			"user_id", oauthcode.UserID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query user")
	}

	var org *entity.Organization
	var orgID *uuid.UUID
	if oauthcode.OrganizationID != nil {
		o, err := uc.orgRepo.GetByID(ctx, *oauthcode.OrganizationID)
		if err != nil {
			uc.logger.Error(ctx, "failed to query organization",
				"action", "EXCHANGE",
				"organization_id", oauthcode.OrganizationID,
				"error", err.Error(),
			)
			return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to query organization")
		}
		org = o
		orgID = &org.ID
	}

	roles, err := uc.userRoleRepo.ListRolesByUser(ctx, user.ID, orgID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query roles",
			"action", "EXCHANGE",
			"user_id", user.ID,
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
			Issuer:    uc.issuerURL,
			Subject:   user.ID.String(),
			Audience:  jwt.ClaimStrings{app.ClientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.Must(uuid.NewV7()).String(),
		},
		Scopes:      oauthcode.Scopes,
		Roles:       roleSlugs,
		Permissions: permSlugs,
	}

	var idToken string
	if slices.Contains(oauthcode.Scopes, "openid") {
		idClaims := &entity.IDTokenClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    uc.issuerURL,
				Subject:   user.ID.String(),
				Audience:  jwt.ClaimStrings{app.ClientID},
				ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
				IssuedAt:  jwt.NewNumericDate(now),
				ID:        uuid.Must(uuid.NewV7()).String(),
			},
			AuthTime: now.Unix(),
		}

		if slices.Contains(oauthcode.Scopes, "profile") {
			idClaims.Name = user.Name
		}

		if slices.Contains(oauthcode.Scopes, "email") {
			idClaims.Email = user.Email
			idClaims.EmailVerified = user.EmailVerifiedAt != nil
		}

		idToken, err = uc.jwtSvc.GenerateIDToken(ctx, idClaims)
		if err != nil {
			uc.logger.Error(ctx, "failed to generate id token",
				"action", "EXCHANGE",
				"user_id", user.ID,
				"error", err.Error())
			return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate id token")
		}
	}

	if org != nil {
		orgIDStr := org.ID.String()
		claims.OrgID = &orgIDStr
		claims.OrgSlug = &org.Slug
	}

	accessToken, err := uc.jwtSvc.GenerateAccessToken(ctx, claims)
	if err != nil {
		uc.logger.Error(ctx, "failed to generate access token",
			"action", "EXCHANGE",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate access token")
	}

	refreshToken, err := uc.jwtSvc.GenerateRefreshToken()
	if err != nil {
		uc.logger.Error(ctx, "failed to generate refresh token",
			"action", "EXCHANGE",
			"user_id", user.ID,
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to generate refresh token")
	}

	refresh := &entity.OAuthRefreshToken{
		ID:             uuid.Must(uuid.NewV7()),
		TokenHash:      hash.HashSHA256(refreshToken),
		UserID:         user.ID,
		ApplicationID:  app.ID,
		OrganizationID: &org.ID,
		Scopes:         oauthcode.Scopes,
		CreatedAt:      now,
		ExpiresAt:      now.AddDate(0, 0, 30),
	}

	if err := uc.oauthRefRepo.Create(ctx, refresh); err != nil {
		uc.logger.Error(ctx, "failed to persist refresh token",
			"action", "EXCHANGE",
			"error", err.Error(),
		)
		return nil, apperr.NewDirectError(enum.SERVER_ERROR, "failed to persist refresh token")
	}

	token := &ExchangeResult{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Scope:        strings.Join(oauthcode.Scopes, " "),
		IDToken:      idToken,
		ExpiresIn:    int(expiresIn.Seconds()),
	}

	return token, nil
}
