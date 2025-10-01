package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/tuusuario/signupgo/internal/di"
	"github.com/tuusuario/signupgo/internal/rest/server"
)

func NewRouter(inj *di.Injector) *gin.Engine {
	r := server.NewEngine()

	// Public
	r.GET("/signup", server.GetSignupList(inj))

	// Protected
	auth := r.Group("")
	auth.Use(server.IntrospectionMiddleware(inj.Cfg.AuthBaseURL, inj.Cfg.HTTPTimeoutMs))
	auth.POST("/signup", server.PostSignup(inj))
	auth.DELETE("/signup/:id", server.DeleteSignup(inj))

	return r
}
