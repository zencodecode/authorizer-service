package auth

import (
	"context"
	"time"

	"github.com/k0kubun/pp"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/application"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/oauthauthorizationcode"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	TokenParams struct {
		GrantType    string
		Code         string
		RedirectURI  string
		ClientID     string
		ClientSecret string
		CodeVerifier string
	}
	TokenResult struct {
		User *entity.User
	}
)

type tokenUsecase struct {
	userRepo      user.Repository
	appRepo       application.Repository
	oauthCodeRepo oauthauthorizationcode.Repository
	roleRepo      role.Repository
	userRoleRepo  userrole.Repository
	permRepo      permission.Repository
	rolePermRepo  rolepermission.Repository
	token         service.JWTService
	logger        service.Logger
}

func NewTokenUsecase(
	userRepo user.Repository,
	appRepo application.Repository,
	oauthCodeRepo oauthauthorizationcode.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
	permRepo permission.Repository,
	rolePermRepo rolepermission.Repository,
	token service.JWTService,
	logger service.Logger,
) TokenUsecase {
	return &tokenUsecase{
		userRepo:      userRepo,
		appRepo:       appRepo,
		oauthCodeRepo: oauthCodeRepo,
		roleRepo:      roleRepo,
		userRoleRepo:  userRoleRepo,
		permRepo:      permRepo,
		rolePermRepo:  rolePermRepo,
		token:         token,
		logger:        logger,
	}
}

func (uc *tokenUsecase) Execute(ctx context.Context, params TokenParams) (*TokenResult, error) {
	app, err := uc.appRepo.GetByClientID(ctx, params.ClientID)
	if err != nil {
		uc.logger.Error(ctx, "failed to query applicaton",
			"action", "TOKEN",
			"client_id", &params.ClientID,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query applicaton")
	}

	if !hash.CheckHash(app.ClientSecretHash, params.ClientSecret) {
		return nil, newAuthError("bad_request", "client secret is invalid")
	}

	codeHash := hash.HashSHA256(params.Code)

	oauthcode, err := uc.oauthCodeRepo.GetByCodeHash(ctx, codeHash)
	if err != nil {
		uc.logger.Error(ctx, "failed to query oauth authorization code",
			"action", "TOKEN",
			"code_hash", codeHash,
			"error", err.Error(),
		)
		return nil, newAuthError("server_error", "failed to query oauth authorization code")
	}

	if !oauthcode.ExpiresAt.After(time.Now()) {
		return nil, newAuthError("not_found", "oauth authorization code has been expired")
	}

	if oauthcode.UsedAt != nil {
		return nil, newAuthError("not_found", "oauth authorization code has been used")
	}

	if oauthcode.RedirectURI != params.RedirectURI {
		return nil, newAuthError("bad_request", "redirect URI does not match")
	}

	pp.Println(app)
	return nil, nil
}
