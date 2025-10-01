// internal/http/routes.go
package http

import (
	"github.com/gin-gonic/gin"

	"github.com/tuusuario/go-auth/internal/config"
	"github.com/tuusuario/go-auth/internal/http/handler"
	"github.com/tuusuario/go-auth/internal/http/middleware"
	"github.com/tuusuario/go-auth/internal/repository"
	"github.com/tuusuario/go-auth/internal/service"
)

// NewRouter construye el router de Gin con todos los endpoints de Auth.
func NewRouter(
	cfg *config.Config,
	userSvc service.UserService,
	sessionRepo repository.SessionRepository,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	uh := handler.NewUserHandler(userSvc)

	// --- Públicos ---
	r.POST("/user", uh.Register)
	r.POST("/user/login", uh.Login)

	// Introspección S2S (para Event/Signup) -> invalida token si la sesión (sid) no existe
	r.GET("/token/introspect", handler.TokenIntrospectHandler(cfg, sessionRepo))

	// --- Protegidos por JWT local (middleware debe setear "username" y "sid") ---
	auth := r.Group("")
	auth.Use(middleware.JWT(cfg.JWTSecret, sessionRepo)) // valida firma/exp y pone username/sid en el contexto
	{
		auth.GET("/user/current", uh.Current)
		auth.POST("/user/logout", uh.Logout)
	}

	return r
}
