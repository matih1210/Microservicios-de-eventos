package handler

import (
    "encoding/base64"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "github.com/tuusuario/go-auth/internal/service"
)

type UserHandler struct{ svc service.UserService }
func NewUserHandler(s service.UserService) *UserHandler { return &UserHandler{svc: s} }

type registerReq struct {
    Name     string `json:"name" binding:"required,min=2"`
    Username string `json:"username" binding:"required,alphanum,min=3"`
    Password string `json:"password" binding:"required,min=6"`
}

func (h *UserHandler) Register(c *gin.Context) {
    var req registerReq
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{}) // 400 si falta algo
        return
    }
    u, err := h.svc.Register(req.Name, req.Username, req.Password)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{})
        return
    }
    c.JSON(http.StatusCreated, u)
}

func (h *UserHandler) Current(c *gin.Context) {
    username := c.GetString("username")
    u, err := h.svc.CurrentByUsername(username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    c.JSON(http.StatusOK, u)
}

// Authorization: Basic base64("user:password")
func (h *UserHandler) Login(c *gin.Context) {
    auth := c.GetHeader("Authorization")
    if !strings.HasPrefix(auth, "Basic ") {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    dec, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(auth, "Basic "))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    parts := strings.SplitN(string(dec), ":", 2)
    if len(parts) != 2 {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    token, err := h.svc.LoginBasic(parts[0], parts[1])
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) Logout(c *gin.Context) {
    sid := c.GetString("sid")
    if sid == "" {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    if err := h.svc.LogoutBySID(sid); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{})
        return
    }
    c.Status(http.StatusOK)
}
