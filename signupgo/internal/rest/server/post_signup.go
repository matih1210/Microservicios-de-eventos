package server

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/tuusuario/signupgo/internal/di"
    "github.com/tuusuario/signupgo/internal/domain/signup"
)

type postSignupReq struct {
    EventID string `json:"eventId" binding:"required"`
}

func PostSignup(inj *di.Injector) gin.HandlerFunc {
    return func(c *gin.Context) {
        var req postSignupReq
        if err := c.ShouldBindJSON(&req); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{})
            return
        }
        userID := c.GetString("uid")
        userName := c.GetString("usr")
        s, err := inj.SignupSvc.Create(c.Request.Context(), userID, userName, req.EventID)
        if err != nil {
            switch err {
            case signup.ErrEventNotFound, signup.ErrEventCanceled:
                c.JSON(http.StatusBadRequest, gin.H{})
            case signup.ErrAlreadySigned:
                c.JSON(http.StatusBadRequest, gin.H{})
            default:
                c.JSON(http.StatusInternalServerError, gin.H{})
            }
            return
        }
        c.JSON(http.StatusCreated, s)
    }
}
