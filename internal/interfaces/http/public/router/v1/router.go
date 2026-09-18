package oauth

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	oidc.GET("/.well-known/openid-configuration", r.container.DiscoveryHandler.OpenIDConfiguration)
	oidc.GET("/.well-known/jwks.json", r.container.DiscoveryHandler.JWKS)
}

func (r *Router) MountOAuth2() {
	oauth2 := r.group.Group("/oauth2")
	oauth2.GET("/authorize", r.container.AuthHandler.Authorize)

	loginCfg := middleware.DefaultRateLimiterConfig()
	loginCfg.MaxRequests = 5
	loginCfg.Window = 5 * time.Minute
	loginCfg.Prefix = "ratelimit:login"
	loginCfg.FailOpen = false

	tokenCfg := middleware.DefaultRateLimiterConfig()
	tokenCfg.MaxRequests = 20
	tokenCfg.Window = time.Minute
	tokenCfg.Prefix = "ratelimit:token"
	tokenCfg.FailOpen = true
	tokenCfg.KeyFunc = clientIDKeyFunc

	forgotCfg := middleware.DefaultRateLimiterConfig()
	forgotCfg.MaxRequests = 3
	forgotCfg.Window = time.Hour
	forgotCfg.Prefix = "ratelimit:forgot-password"
	forgotCfg.FailOpen = false
	forgotCfg.KeyFunc = emailBodyKeyFunc

	registerCfg := middleware.DefaultRateLimiterConfig()
	registerCfg.MaxRequests = 5
	registerCfg.Window = time.Hour
	registerCfg.Prefix = "ratelimit:register"
	registerCfg.FailOpen = false

	oauth2.POST("/authorize/login", middleware.RateLimiter(r.container.ConnManager.Redis.GetClient(), loginCfg), r.container.AuthHandler.Login)
	oauth2.POST("/authorize/consent", r.container.AuthHandler.Consent)

	oauth2.POST("/token", middleware.RateLimiter(r.container.ConnManager.Redis.GetClient(), tokenCfg), r.container.TokenHandler.Token)
	oauth2.POST("/revoke", r.container.TokenHandler.Revoke)

	oauth2.POST("/register", middleware.RateLimiter(r.container.ConnManager.Redis.GetClient(), registerCfg), r.container.AuthHandler.Register)
	oauth2.GET("/verify-email", r.container.AuthHandler.VerifyEmail)

	oauth2.POST("/password/forgot", middleware.RateLimiter(r.container.ConnManager.Redis.GetClient(), forgotCfg), r.container.AuthHandler.ForgotPassword)
	oauth2.POST("/password/reset", r.container.AuthHandler.ResetPassword)

	oauth2.GET("/userinfo", middleware.JWTAuth(r.container.JWTService, r.container.Logger), r.container.AuthHandler.UserInfo)
	oauth2.POST("/logout", middleware.JWTAuth(r.container.JWTService, r.container.Logger), r.container.AuthHandler.Logout)
}

func clientIDKeyFunc(c *gin.Context) string {
	if claims, _ := middleware.GetClaims(c); claims != nil {
		var audString string
		if len(claims.Audience) > 0 {
			audString = claims.Audience[0]
		}
		return "client:" + audString
	}

	if id := c.PostForm("client_id"); id != "" {
		return "client:" + id
	}

	return "unknown-client:" + c.ClientIP()
}

type forgotPasswordBody struct {
	Email string `json:"email" binding:"required,email"`
}

func emailBodyKeyFunc(c *gin.Context) string {
	var body forgotPasswordBody
	if err := c.ShouldBindBodyWith(&body, binding.JSON); err != nil || body.Email == "" {
		return "invalid-email:" + c.ClientIP()
	}
	return "email:" + strings.ToLower(strings.TrimSpace(body.Email))
}
