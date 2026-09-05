package auth

import (
	"github.com/zencodecode/authorizer-service/internal/domain/service"
	"github.com/zencodecode/authorizer-service/internal/usecase/auth"
)

type Handler struct {
	authorizeUC   auth.AuthorizeUsecase
	loginUC       auth.LoginUsecase
	consentUC     auth.ConsentUsecase
	logoutUC      auth.LogoutUsecase
	userInfoUC    auth.UserInfoUsecase
	registerUC    auth.RegisterUsecase
	verifyEmailUC auth.VerifyEmailUsecase
	jwtSvc        service.JWTService
	logger        service.Logger
}

func New(
	authorizeUC auth.AuthorizeUsecase,
	loginUC auth.LoginUsecase,
	consentUC auth.ConsentUsecase,
	logoutUC auth.LogoutUsecase,
	userInfoUC auth.UserInfoUsecase,
	registerUC auth.RegisterUsecase,
	verifyEmailUC auth.VerifyEmailUsecase,
	jwtSvc service.JWTService,
	logger service.Logger,
) *Handler {
	return &Handler{
		authorizeUC:   authorizeUC,
		loginUC:       loginUC,
		consentUC:     consentUC,
		logoutUC:      logoutUC,
		userInfoUC:    userInfoUC,
		registerUC:    registerUC,
		verifyEmailUC: verifyEmailUC,
		jwtSvc:        jwtSvc,
		logger:        logger,
	}
}
