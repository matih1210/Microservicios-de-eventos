package server

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/tuusuario/eventgo/internal/di"
    "github.com/tuusuario/eventgo/internal/domain/event"
)

func GetEventByID(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        e, err := inj.EventSvc.GetByID(c.Request.Context(), id)
        if err != nil {
            if err == event.ErrNotFound {
                c.JSON(http.StatusNotFound, gin.H{})
                return
            }
            c.JSON(http.StatusInternalServerError, gin.H{})
            return
        }
        c.JSON(http.StatusOK, gin.H{
            "id": e.ID,
            "name": e.Name,
            "when": e.When,
            "updated": e.Updated,
            "created": e.Created,
            "canceled": e.Canceled,
            "ownerId": e.OwnerID,
        })
    }
}
