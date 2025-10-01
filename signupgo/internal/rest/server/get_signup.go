package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/tuusuario/signupgo/internal/di"
)

func GetSignupList(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        eventID := c.Query("eventId")
        if eventID == "" {
            c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        items, err := inj.SignupSvc.ListByEvent(c.Request.Context(), eventID)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{})
            return
        }
        c.JSON(http.StatusOK, items)
    }
}
