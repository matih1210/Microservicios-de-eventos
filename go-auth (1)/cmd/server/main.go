package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/tuusuario/go-auth/internal/config"
	userhandler "github.com/tuusuario/go-auth/internal/http/handler"
	"github.com/tuusuario/go-auth/internal/http/middleware"
	mongorepo "github.com/tuusuario/go-auth/internal/repository/mongo"
	"github.com/tuusuario/go-auth/internal/service"
)

func main() {
	_ = godotenv.Load()

	cfg := config.New()

	mongoClient, err := mongorepo.NewClient(cfg.MongoURI)
	if err != nil {
		log.Fatal("mongo connect error: ", err)
	}
	defer mongoClient.Disconnect()

	userRepo := mongorepo.NewUserRepository(mongoClient, cfg.MongoDB)
	sessionRepo := mongorepo.NewSessionRepository(mongoClient, cfg.MongoDB)
	userSvc := service.NewUserService(userRepo, sessionRepo, cfg.JWTSecret, cfg.JWTExpMin)

	r := gin.Default()

	uh := userhandler.NewUserHandler(userSvc)

	// --- Públicos ---
	r.POST("/user", uh.Register)
	r.POST("/user/login", uh.Login)

	// Introspección S2S (Event/Signup llaman acá para validar sesión/Logout)
	r.GET("/token/introspect", userhandler.TokenIntrospectHandler(cfg, sessionRepo))

	// --- Protegidos (JWT local; el middleware debe setear username y sid) ---
	auth := r.Group("")
	auth.Use(middleware.JWT(cfg.JWTSecret, sessionRepo))
	auth.GET("/user/current", uh.Current)
	auth.POST("/user/logout", uh.Logout)

	addr := ":" + cfg.Port
	log.Println("Auth server listening on", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
