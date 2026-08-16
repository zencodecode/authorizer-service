package auth

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zencodecode/authorizer-service/internal/domain/entity"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/permission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/role"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/rolepermission"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/user"
	"github.com/zencodecode/authorizer-service/internal/domain/repository/userrole"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/pkg/hash"
)

type (
	LoginParams struct {
		Email    string
		Password string
		OrgID    string
	}

	LoginOutput struct {
		User         *entity.User
		Token        string
		RefreshToken string
		Claims       *entity.Claims
	}
)

type loginUsecase struct {
	userRepo     user.Repository
	roleRepo     role.Repository
	userRoleRepo userrole.Repository
	permRepo     permission.Repository
	rolePermRepo rolepermission.Repository
	token        service.Token
	logger       service.Logger
}

func NewLoginUsecase(
	userRepo user.Repository,
	roleRepo role.Repository,
	userRoleRepo userrole.Repository,
	permRepo permission.Repository,
	rolePermRepo rolepermission.Repository,
	token service.Token,
	logger service.Logger,
) LoginUsecase {
	return &loginUsecase{
		userRepo:     userRepo,
		roleRepo:     roleRepo,
		userRoleRepo: userRoleRepo,
		permRepo:     permRepo,
		rolePermRepo: rolePermRepo,
		token:        token,
		logger:       logger,
	}
}

func (uc *loginUsecase) Execute(ctx context.Context, params LoginParams) (*LoginOutput, error) {
	// var authorization []entity.Authorization
	var audiences []string

	user, err := uc.userRepo.GetByEmail(ctx, params.Email)
	if err != nil {
		uc.logger.Warn(ctx, "failed to fetch user by email ",
			"email", params.Email,
			"context", "LOGIN",
			"error", err.Error(),
		)
		return nil, errors.New("email or password is invalid")
	}

	if !hash.CheckHash(user.PasswordHash, params.Password) {
		uc.logger.Warn(ctx, "invalid password",
			"user_id", user.ID,
			"context", "LOGIN",
		)
		return nil, errors.New("email or password is invalid")
	}

	roles, err := uc.userRoleRepo.ListRolesByUser(ctx, user.ID.String(), &params.OrgID)
	if len(roles) > 0 {
		roleSet := make(map[string]struct{})
		permSet := make(map[string]struct{})

		for _, r := range roles {
			roleSet[r.Slug] = struct{}{}

			perms, _ := uc.rolePermRepo.ListPermissionsByRole(ctx, r.ID.String())
			for _, p := range perms {
				permSet[p.Slug] = struct{}{}
			}
		}

		// authorization = append(authorization, entity.Authorization{
		// 	Roles:       mapKeys(roleSet),
		// 	Permissions: mapKeys(permSet),
		// })

		audiences = append(audiences, "LOG-SERVICE")
	}

	now := time.Now()
	claims := &entity.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   user.ID.String(),
			Audience:  audiences,
			ExpiresAt: now.Add(time.Hour).Unix(),
		},
	}

	accessToken, err := uc.token.GenerateAccessToken(ctx, claims)
	if err != nil {
		uc.logger.Error("Failed to generate access token", service.Fields{
			"user_id": user.ID,
			"context": "LOGIN",
			"error":   err.Error(),
		})
		return nil, errors.New("failed to generate access token")
	}

	refreshToken, err := uc.token.GenerateRefreshToken()
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", service.Fields{
			"user_id": user.ID,
			"context": "LOGIN",
			"error":   err.Error(),
		})
		return nil, errors.New("failed generating refresh token")
	}

	err = uc.token.StoreRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		uc.logger.Error("Failed to store refresh token", service.Fields{
			"user_id": user.ID,
			"context": "LOGIN",
			"error":   err.Error(),
		})
		return nil, errors.New("failed saving refresh token")
	}

	output := &LoginOutput{
		User:         user,
		Token:        accessToken,
		RefreshToken: refreshToken,
		Claims:       claims,
	}
	return output, nil
}

func mapKeys(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
