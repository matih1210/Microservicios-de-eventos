package server

import "github.com/gin-gonic/gin"

func NewEngine() *gin.Engine {
    r := gin.New()
    r.Use(gin.Logger(), gin.Recovery())
    return r
}
