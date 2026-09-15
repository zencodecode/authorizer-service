package bootstrap

import (
	"github.com/redis/go-redis/v9"
	"github.com/zencodecode/authorizer-service/internal/config"
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/auth"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/database"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/email"
	"github.com/zencodecode/authorizer-service/internal/infrastructure/driver/rabbitmq"
	"gorm.io/gorm"

	// PostgreSQL repositories
	pgRepo "github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/postgres/repository"
	// Redis repositories
	redisRepo "github.com/zencodecode/authorizer-service/internal/infrastructure/persistence/redis/repository"
	// Usecases
	// auditSvc "github.com/zencodecode/authorizer-service/internal/usecase/audit"
	authUC "github.com/zencodecode/authorizer-service/internal/usecase/auth"
	tokenUC "github.com/zencodecode/authorizer-service/internal/usecase/token"
	userUC "github.com/zencodecode/authorizer-service/internal/usecase/user"

	// Handlers
	authHandler "github.com/zencodecode/authorizer-service/internal/interfaces/http/public/handler/v1/auth"
	discoveryHandler "github.com/zencodecode/authorizer-service/internal/interfaces/http/public/handler/v1/discovery"
	tokenHandler "github.com/zencodecode/authorizer-service/internal/interfaces/http/public/handler/v1/token"
	// Private handlers
	// userPrivHandler "github.com/zencodecode/authorizer-service/internal/interfaces/httpprivate/handler/v1/user"
)

type Container struct {
	Config      config.Config
	ConnManager *database.ConnectionManager
	Amqp        *rabbitmq.Connection
	Logger      service.Logger
	// Handlers — Public
	DiscoveryHandler *discoveryHandler.Handler
	AuthHandler      *authHandler.Handler
	TokenHandler     *tokenHandler.Handler
	// Handlers — Private
	// UserHandler *userPrivHandler.Handler
	// OrgHandler      *orgPrivHandler.Handler
	// AppHandler      *appPrivHandler.Handler
	// RoleHandler     *rolePrivHandler.Handler
	// PermHandler     *permPrivHandler.Handler
	// InvHandler      *invPrivHandler.Handler
	// Services
	JWTService service.JWTService
	// AuditSvc   auditSvc.Service
}

func NewContainer(cfg config.Config, cm *database.ConnectionManager, amqp *rabbitmq.Connection, logger service.Logger) *Container {
	return &Container{
		Config:      cfg,
		ConnManager: cm,
		Amqp:        amqp,
		Logger:      logger,
	}
}
func (c *Container) Build(db *gorm.DB, redisClient *redis.Client) {
	// ──────── Config ────────
	issuer := c.Config.Auth.OIDC.Issuer
	authCodeExpiry := c.Config.Auth.OIDC.AuthorizationCodeExpiry
	privateKey := c.Config.Auth.JWT.PrivateKey
	publicKey := c.Config.Auth.JWT.PublicKey
	keyID := c.Config.Auth.JWT.KeyID
	tokenExpiry := c.Config.Auth.JWT.TokenExpiry

	// ──────── Repositories (PostgreSQL) ────────
	userRepo := pgRepo.NewUserRepository(db)
	orgRepo := pgRepo.NewOrganizationRepository(db)
	orgUserRepo := pgRepo.NewOrganizationUserRepository(db)
	orgAppRepo := pgRepo.NewOrganizationApplicationRepository(db)
	appRepo := pgRepo.NewApplicationRepository(db)
	scopeRepo := pgRepo.NewApplicationScopeRepository(db)
	roleRepo := pgRepo.NewRoleRepository(db)
	permRepo := pgRepo.NewPermissionRepository(db)
	rolePermRepo := pgRepo.NewRolePermissionRepository(db)
	userRoleRepo := pgRepo.NewUserRoleRepository(db)
	auditLogRepo := pgRepo.NewAuditLogRepository(db)
	verifRepo := pgRepo.NewEmailVerificationTokenRepository(db)
	resetRepo := pgRepo.NewPasswordResetTokenRepository(db)
	invRepo := pgRepo.NewInvitationRepository(db)

	// ──────── Repositories (Redis) ────────
	sessionRepo := redisRepo.NewAuthorizeSessionRepository(redisClient)
	oauthCodeRepo := redisRepo.NewOAuthAuthorizationCodeRepository(redisClient)
	oauthRefRepo := redisRepo.NewOAuthRefreshTokenRepository(redisClient)

	// ──────── Infrastructure Services ────────
	jwtSvc := auth.NewJWTService(
		privateKey,
		keyID,
		tokenExpiry,
	)
	c.JWTService = jwtSvc
	authCodeSvc := auth.NewAuthorizationCode(oauthCodeRepo, sessionRepo)

	// ──────── Audit Service ────────
	// c.AuditSvc = auditSvc.NewAuditService(auditLogRepo, c.Logger)

	// ──────── Email Publisher ────────
	emailChan := c.Amqp.Channel()
	emailPublisher, _ := email.NewPublisher(emailChan, c.Logger)

	// ──────── Usecases — Auth (Public) ────────
	authorizeUC := authUC.NewAuthorizeUsecase(appRepo, scopeRepo, sessionRepo, c.Logger)
	loginUC := authUC.NewLoginUsecase(userRepo, sessionRepo, appRepo, orgUserRepo, orgRepo,
		auditLogRepo, authCodeSvc, authCodeExpiry, c.Logger)
	consentUC := authUC.NewConsentUsecase(userRepo, orgUserRepo, orgAppRepo, orgRepo, appRepo,
		sessionRepo, authCodeSvc, c.Logger)
	logoutUC := authUC.NewLogoutUsecase(userRepo, appRepo, oauthRefRepo, c.Logger)
	userInfoUC := authUC.NewUserInfoUsecase(userRepo, c.Logger)
	registerUC := authUC.NewRegisterUsecase(userRepo, verifRepo, emailPublisher, &c.Config, c.Logger)
	verifyEmailUC := authUC.NewVerifyEmailUsecase(verifRepo, userRepo, c.Logger)
	forgotPwdUC := authUC.NewForgotPasswordUsecase(userRepo, resetRepo, emailPublisher, &c.Config, c.Logger)
	resetPwdUC := authUC.NewResetPasswordUsecase(userRepo, resetRepo, oauthRefRepo, c.Logger)

	// ──────── Usecases — Token (Public) ────────
	exchangeUC := tokenUC.NewExchangeUsecase(userRepo, appRepo, oauthCodeRepo, orgRepo,
		roleRepo, userRoleRepo, permRepo, rolePermRepo, oauthRefRepo, jwtSvc, issuer, c.Logger)
	refreshUC := tokenUC.NewRefreshUsecase(userRepo, appRepo, orgRepo, userRoleRepo,
		rolePermRepo, permRepo, oauthRefRepo, jwtSvc, issuer, c.Logger)
	revokeUC := tokenUC.NewRevokeUsecase(appRepo, oauthRefRepo, c.Logger)

	// ──────── Usecases — Admin (Private) ────────
	_ = userUC.NewCreateUsecase // TODO: wire admin usecases

	// ──────── Handlers — Public ────────
	c.DiscoveryHandler = discoveryHandler.New(
		publicKey,
		keyID,
		issuer,
	)
	c.AuthHandler = authHandler.New(authorizeUC, loginUC, consentUC, logoutUC, userInfoUC, registerUC,
		verifyEmailUC, forgotPwdUC, resetPwdUC, jwtSvc, c.Logger)
	c.TokenHandler = tokenHandler.New(exchangeUC, refreshUC, revokeUC, c.Logger)
	// ──────── Handlers — Private ────────
	// c.UserHandler = userPrivHandler.NewHandler(userUC, c.Logger)
	// c.OrgHandler = ...
	// c.AppHandler = ...
	// c.RoleHandler = ...
	// c.PermHandler = ...
	// Ignore unused
	_, _, _, _, _, _ = resetRepo, invRepo, orgAppRepo, roleRepo, permRepo, rolePermRepo
}
