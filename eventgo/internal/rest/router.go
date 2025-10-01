package rest

import (
	"github.com/gin-gonic/gin"

	"github.com/tuusuario/eventgo/internal/di"
	"github.com/tuusuario/eventgo/internal/rest/server"
)

func NewRouter(inj *di.Injector) *gin.Engine {
	r := server.NewEngine()

	// Public routes
	r.GET("/event", server.GetEventList(inj))
	r.GET("/event/:id", server.GetEventByID(inj))

	// Auth protected
	auth := r.Group("")
	auth.Use(server.IntrospectionMiddleware(inj.Cfg.AuthBaseURL, inj.Cfg.HTTPTimeoutMs))

	auth.POST("/event", server.PostEvent(inj))
	auth.PUT("/event/:id", server.PutEvent(inj))
	auth.DELETE("/event/:id", server.DeleteEvent(inj))

	return r
}
