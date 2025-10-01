package server

import (
    "net/http"

    "github.com/gin-gonic/gin"

    "github.com/tuusuario/eventgo/internal/di"
    "github.com/tuusuario/eventgo/internal/domain/event"
)

type putEventReq struct {
    ID   string `json:"id" binding:"required"`
    Name string `json:"name" binding:"required,min=2"`
    When int64  `json:"when" binding:"required"`
}

func PutEvent(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req putEventReq
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        if req.ID != c.Param("id") {
            c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        actorID := c.GetString("uid")
        e, err := inj.EventSvc.Update(c.Request.Context(), req.ID, req.Name, req.When, actorID)
        if err != nil {
            switch err {
            case event.ErrPastDate:
                c.JSON(http.StatusBadRequest, gin.H{})
            case event.ErrForbidden:
                c.JSON(http.StatusForbidden, gin.H{})
            case event.ErrNotFound:
                c.JSON(http.StatusNotFound, gin.H{})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{})
            }
            return
        }
        c.JSON(http.StatusOK, e)
    }
}
