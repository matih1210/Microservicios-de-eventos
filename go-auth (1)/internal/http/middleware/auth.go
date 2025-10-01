package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tuusuario/go-auth/internal/repository"
)

func JWT(secret string, sessions repository.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		tok, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !tok.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims := tok.Claims.(jwt.MapClaims)
		usr, _ := claims["usr"].(string)
		sid, _ := claims["sid"].(string)
		if usr == "" || sid == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// validar que la sesión siga existiendo (logout => 401)
		if sessions != nil {
			ok, err := sessions.Exists(c.Request.Context(), sid)
			if err != nil || !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
		}

		c.Set("username", usr)
		c.Set("sid", sid)
		c.Next()
	}
}
