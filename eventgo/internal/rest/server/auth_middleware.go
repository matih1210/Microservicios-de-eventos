package server

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type introspectOK struct {
	Active bool   `json:"active"`
	UID    string `json:"uid"`
	USR    string `json:"usr"`
	SID    string `json:"sid"`
	Exp    int64  `json:"exp"`
}

func IntrospectionMiddleware(authBase string, timeoutMs int) gin.HandlerFunc {
	if timeoutMs <= 0 {
		timeoutMs = 2000
	}
	client := &http.Client{Timeout: time.Duration(timeoutMs) * time.Millisecond}

	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}
		req, _ := http.NewRequestWithContext(c.Request.Context(), "GET", authBase+"/token/introspect", nil)
		req.Header.Set("Authorization", auth)

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			if resp != nil {
				_ = resp.Body.Close()
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}
		defer resp.Body.Close()

		var body introspectOK
		if err := json.NewDecoder(resp.Body).Decode(&body); err != nil || !body.Active || body.UID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{})
			return
		}
		// pasar identidad a los handlers
		c.Set("uid", body.UID)
		c.Set("usr", body.USR)
		c.Set("sid", body.SID)

		c.Next()
	}
}
