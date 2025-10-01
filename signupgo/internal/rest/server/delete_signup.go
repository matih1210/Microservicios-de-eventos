package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/tuusuario/signupgo/internal/di"
    "github.com/tuusuario/signupgo/internal/domain/signup"
)

func DeleteSignup(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        id := c.Param("id")
        actorID := c.GetString("uid")
        if err := inj.SignupSvc.Cancel(c.Request.Context(), id, actorID); err != nil {
            switch err {
            case signup.ErrForbidden:
                c.JSON(http.StatusForbidden, gin.H{})
            case signup.ErrNotFound:
                c.JSON(http.StatusNotFound, gin.H{})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{})
            }
            return
        }
        c.Status(http.StatusOK)
    }
}
