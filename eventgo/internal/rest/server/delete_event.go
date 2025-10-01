package server

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/tuusuario/eventgo/internal/di"
    "github.com/tuusuario/eventgo/internal/domain/event"
)

func DeleteEvent(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        actorID := c.GetString("uid")
        if err := inj.EventSvc.Cancel(c.Request.Context(), id, actorID); err != nil {
            switch err {
            case event.ErrForbidden:
                c.JSON(http.StatusForbidden, gin.H{})
            case event.ErrNotFound:
                c.JSON(http.StatusNotFound, gin.H{})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{})
            }
            return
        }
        c.Status(http.StatusOK)
    }
}
