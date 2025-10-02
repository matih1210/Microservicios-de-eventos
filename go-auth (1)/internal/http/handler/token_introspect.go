package handler

// El handler de introspección es lo que permite que otros microservicios (Event, Signup) no confíen solo en validar localmente el JWT, sino que le pregunten a Auth si el token todavía es válido y la sesión sigue viva.
// Así es como lograste que el logout invalide tokens automáticamente 🚀.

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/tuusuario/go-auth/internal/config"
	"github.com/tuusuario/go-auth/internal/repository"
)

type IntrospectResp struct {
	Active bool   `json:"active"`
	UID    string `json:"uid"`
	USR    string `json:"usr"`
	SID    string `json:"sid"`
	Exp    int64  `json:"exp"`
}

func TokenIntrospectHandler(cfg *config.Config, sessions repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		// 1) validar firma/exp
		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := tok.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		uid, _ := claims["uid"].(string)
		usr, _ := claims["usr"].(string)
		sid, _ := claims["sid"].(string)
		expF, _ := claims["exp"].(float64)
		if uid == "" || sid == "" || time.Now().Unix() >= int64(expF) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 2) validar sesión aún existe (logout la borra)
		okExists, err := sessions.Exists(c.Request.Context(), sid)
		if err != nil || !okExists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.JSON(http.StatusOK, IntrospectResp{
			Active: true, UID: uid, USR: usr, SID: sid, Exp: int64(expF),
		})
	}
}
