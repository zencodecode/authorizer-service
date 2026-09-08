package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/zencodecode/authorizer-service/internal/bootstrap"
	"github.com/zencodecode/authorizer-service/internal/interfaces/http/middleware"
)

type Router struct {
	container *bootstrap.Container
	group     *gin.RouterGroup
}

func SetupRouter(c *bootstrap.Container, group *gin.RouterGroup) *Router {
	return &Router{container: c, group: group}
}

func (r *Router) MountOIDC() {
	oidc := r.group.Group("")
	oidc.POST("/.well-known/openid-configuration", r.container.DiscoveryHandler.OpenIDConfiguration)
	oidc.POST("/.well-known/jwks.json", r.container.DiscoveryHandler.JWKS)
}

func (r *Router) MountOAuth2() {
	oauth2 := r.group.Group("/oauth2")
	oauth2.GET("/authorize", r.container.AuthHandler.Authorize)
	oauth2.POST("/authorize/login", r.container.AuthHandler.Login)
	oauth2.POST("/authorize/consent", r.container.AuthHandler.Consent)

	oauth2.POST("/token", r.container.TokenHandler.Token)
	oauth2.POST("/revoke", r.container.TokenHandler.Revoke)

	oauth2.POST("/register", r.container.AuthHandler.Register)
	oauth2.GET("/verify-email", r.container.AuthHandler.VerifyEmail)

	oauth2.POST("/password/forgot", r.container.AuthHandler.ForgotPassword)
	oauth2.POST("/password/reset", r.container.AuthHandler.ResetPassword)

	oauth2.GET("/userinfo", middleware.JWTAuth(r.container.JWTService, r.container.Logger), r.container.AuthHandler.UserInfo)
	oauth2.POST("/userinfo", middleware.JWTAuth(r.container.JWTService, r.container.Logger), r.container.AuthHandler.Logout)
}
