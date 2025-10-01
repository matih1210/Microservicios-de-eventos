package server

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"

    "github.com/tuusuario/eventgo/internal/di"
)

type postEventReq struct {
    Name string `json:"name" binding:"required,min=2"`
    When int64  `json:"when" binding:"required"`
}

func PostEvent(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req postEventReq
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        ownerID := c.GetString("uid")
        ownerName := c.GetString("usr")
        e, err := inj.EventSvc.Create(c.Request.Context(), req.Name, req.When, ownerID, ownerName)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "now": time.Now().Unix()})
            return
        }
        c.JSON(http.StatusCreated, e)
    }
}
