package server

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/tuusuario/eventgo/internal/di"
)

func GetEventList(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        items, err := inj.EventSvc.GetOpen(c.Request.Context())
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{})
            return
        }
        c.JSON(http.StatusOK, items)
    }
}
